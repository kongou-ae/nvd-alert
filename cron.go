package main

import (
    "fmt"
    "net"
    "os/exec"
    "github.com/mattn/go-pipeline"
    "strings"
)


func fetchNvd(){
    cmd := exec.Command("go-cve-dictionary","fetchnvd","-last2y")
    cmd.Start()
    fmt.Println("fetching NVD data....")
    cmd.Wait()
    fmt.Println("finished.")
}


func ctrlServer(arg string){
    if arg == "start" {
        cmd := exec.Command("nohup","go-cve-dictionary","server")
        cmd.Start()
        fmt.Println("go-cve-dictionary started")
    }

    if arg == "stop"{
        out, err := pipeline.Output(
            []string{"ps", "-aux"},
            []string{"grep", "go-cve"},
            []string{"grep", "-v", "grep"},
            []string{"awk", "{print $2}"},
        )

        if err != nil {
            fmt.Println(err)
        }
        
        // なぜ改行が入っちゃうんだろう。。。
        pid := strings.TrimRight(string(out), "\n")

        fmt.Println("The pid of go-cve-dictionary is " + pid)
        cmd := exec.Command("kill", pid)
        cmd.Start()
        cmd.Wait()
        fmt.Println("go-cve-dictionary stopped")
    }
}

func main() {
    _, err := net.Listen("tcp", "localhost:1323")
    if err != nil {
        ctrlServer("stop")
    }
    fetchNvd()
    ctrlServer("start")
}