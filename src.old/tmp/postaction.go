package main

import (
  "net/http"
  "text/template"
  "encoding/json"
  "bytes"
  "fmt"
)

func main() {
  command := `
  {"Id":"00000000-0000-0000-0000-000000000000","TargetNode": "{{.Uuid}}", "Command":"run-script", "Args": [ "choco install treesizefree -y" ]}
  `
  jsondata:= `
  { "Uuid": "bc6e31f1-c7a8-11e8-86d2-0a002700000c" }
  `
  t := template.Must(template.New("").Parse(command))

  m := map[string]interface{}{}
  if err := json.Unmarshal([]byte(jsondata), &m); err != nil {
    panic(err)
  }

  var tpl bytes.Buffer
  if err := t.Execute(&tpl, m); err != nil {
    fmt.Printf("Error executing template %s\n", err)
  }
  http.Post("http://controluka.imim.science:8080/job", "application/json", bytes.NewReader(tpl.Bytes()))
}
