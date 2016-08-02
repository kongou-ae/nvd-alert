package main

import (
    "fmt"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/bitly/go-simplejson"
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
    "io/ioutil"
	"strings"
    "net/http"
    "bytes"
//    "strconv"
//    "strings"
//    "encoding/json"
    "html/template"
//    "reflect"
    "time"
    "regexp"
)

func LoadConfig(FilePass string) *simplejson.Json {

    JsonFile, err := ioutil.ReadFile(FilePass)
 	if err != nil { panic(err) }   
    config, err := simplejson.NewJson(JsonFile)
 	if err != nil { panic(err) }   
	return config
}

// get information from sqlite3 by argment
func QuerySqlite3(DbPass string, query string) []string {
    // initialize var
    var cpe_name string
    var dbfile string = DbPass
    result := []string{}

    // connect db
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil { panic(err) }

    // query db	
    rows, err := db.Query(query)
    if err != nil { panic(err) }
    
    // append slice to result
	for rows.Next() {
		err = rows.Scan(&cpe_name)
		result = append(result,cpe_name)
	}
	// make unique slice 
	result = GetUniqueSlice(result)
    return result
}

// Be unique
func GetUniqueSlice(OldSlice []string) []string {
    uniqueSlice := make([]string,0,len(OldSlice))
    encountered := map[string]bool{}

    for i := 0; i < len(OldSlice); i++ {
        if !encountered[OldSlice[i]] {
            encountered[OldSlice[i]] = true
            uniqueSlice = append(uniqueSlice,OldSlice[i])
        }
    }
    return uniqueSlice
}

// make Query by using the target in config.json 
func MakeQuery(target string) (query string){
    targetAry := strings.Split(target,":")
    return "SELECT cpe_name FROM cpes WHERE vendor = '" + targetAry[0] + "' and product = '" + targetAry[1] + "';"
}

func getCveInfobyCpes(cpes string) []interface{}{

    client := &http.Client{}
    var jsonStr = []byte(`{"name": "` + cpes + `"}`)

    // Request を生成
    req, err := http.NewRequest(
        "POST", 
        "http://localhost:1323/cpes",
        bytes.NewBuffer(jsonStr),
    )
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Accept", "application/json")
    if err != nil {
        fmt.Println(err)
    }
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    cves,err := simplejson.NewJson(body)
    if err != nil {
        fmt.Println(err)
    }
    return cves.MustArray()
}

func getCveInfobyCve(cve string) map[string]interface {}{

    client := &http.Client{}

    // Request を生成
    req, err := http.NewRequest(
        "GET", 
        "http://localhost:1323/cves/" + cve,
        nil,
    )
    if err != nil {
        fmt.Println(err)
    }
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    cves,err := simplejson.NewJson(body)
    if err != nil {
        fmt.Println(err)
    }
    m, _ := cves.Map()
    
    return m
}

func sendMailBySendGrid(mailConfig map[string]interface {}, body string,target string)  {
    
    from := mail.NewEmail("", mailConfig["fromAddress"].(string))
    subject := "NVD Alert" 
    to := mail.NewEmail("", mailConfig["toAddress"].(string))
    content := mail.NewContent("text/html", body)

    m := mail.NewV3MailInit(from, subject, to, content)
    request := sendgrid.GetRequest(mailConfig["apikey"].(string), "/v3/mail/send", "https://api.sendgrid.com")
    request.Method = "POST"
    request.Body = mail.GetRequestBody(m)
    
    _, err := sendgrid.API(request)
    if err != nil {
        fmt.Println(err)
    }
    
}

func getHtmlMailBody(cvesInfoDetail []map[string]interface {} ,target string ) string{

    type email struct {
        Target string
        CvesInfoDetail []map[string]interface {}
    }
    
    myEmail := email{}
    myEmail.Target = target
    myEmail.CvesInfoDetail = cvesInfoDetail

    // Create a template using template.html
    tmpl := template.Must(template.ParseFiles("./template.tpl"))
    var buff bytes.Buffer

    // Send the parsed template to buff 
    err := tmpl.Execute(&buff, myEmail)
    if err != nil {
        fmt.Println(err)  
    }
    
    body := buff.String()
    return body    
}

func main() {
    var query string
    var cpes []string
    var cvesInfo []interface{}
    var cvesSlice []string
    var cveInfoDetail map[string]interface {}
    var cvesInfoDetail []map[string]interface {}
    
	config := LoadConfig("./config.json")
    DbPass := config.Get("DbPass").MustString()
    target := config.Get("target").MustArray()
    UpdatePeriod := config.Get("UpdatePeriod").MustInt()
    // forで回すとcpesごとに判断時間が変わってしまうので、ここで定義する
    checkTime := time.Now().Add(-time.Duration(UpdatePeriod) * time.Second)
    
    // コンフィグに記載されているターゲット分の処理を実施
    for i := 0; i < len(target); i++ {
        query = MakeQuery(target[i].(string))
    	cpes = QuerySqlite3(DbPass,query)

        // 取得したcpes分の処理を実施
        cvesSlice = []string{} // スライスの初期化ってこれでいいのか？

        for j :=0; j < len(cpes); j++{
            cvesInfo = getCveInfobyCpes(cpes[j])

            // cpesを使った取得したcveの一覧を作成
            for k := 0; k < len(cvesInfo); k++{
                cvesSlice = append(cvesSlice,cvesInfo[k].(map[string]interface{})["CveID"].(string))
                
            }
            cvesSlice = GetUniqueSlice(cvesSlice)
        }
        
        // 1つのターゲットに紐づくCVEの詳細を作成する
        cvesInfoDetail = make([]map[string]interface {},0) // スライスの初期化ってこれでいいのか？
        for l :=0; l < len(cvesSlice); l++{
            cveInfoDetail = getCveInfobyCve(cvesSlice[l])

            rep := regexp.MustCompile(`.[0-9]+-[0-9]{2}:[0-9]{2}$`)
            updateTime, _ := time.Parse("2006-01-02T15:04:05", rep.ReplaceAllString(cveInfoDetail["Nvd"].(map[string]interface{})["LastModifiedDate"].(string),""))

            if updateTime.After(checkTime){
                cvesInfoDetail = append(cvesInfoDetail,cveInfoDetail)
            }
        }
        
        // cvesInfoDetailに情報が入っていたら
        if len(cvesInfoDetail) != 0 {
            mailBody := getHtmlMailBody(cvesInfoDetail,target[i].(string))
            mailConfig,_ := config.Get("Mail").Map()
            sendMailBySendGrid(mailConfig,mailBody,target[i].(string))            
        }

    }
}