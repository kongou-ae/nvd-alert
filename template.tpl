<h1>NVD Alert</h1>
<h2>{{.Target}}</h2>
<p>This mail was sent by nvd-alert. The new CVEs are as below.</p>
{{ range .CvesInfoDetail }}
    <table border=1 style="margin-bottom:10px">
        <tr>
            <th>ID</th>
            <td>{{.CveID}}</td>
        </tr>
        {{with .Nvd}}
            <tr>
                <th>Score</th>
                <td>{{.Score}}</td>
            </tr>
            <tr>
                <th>Summary</th>
                <td>{{.Summary}}</td>
            </tr>
            <tr>
                <th>LastModifiedDate</th>
                <td>{{.LastModifiedDate}}</td>
            </tr>
        {{end}}
        <tr>
            <th>URL</th>
            <td><a href="https://web.nvd.nist.gov/view/vuln/detail?vulnId={{.CveID}}">https://web.nvd.nist.gov/view/vuln/detail?vulnId={{.CveID}}</a></td>
        </tr>
    </table>
{{ end }}
