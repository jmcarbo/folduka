package main

import (
  "github.com/kataras/iris"
  "fmt"
  "database/sql"
  "github.com/elgs/gosqljson"
  _ "github.com/lib/pq"
  _ "github.com/go-sql-driver/mysql"
  "encoding/json"
  "io/ioutil"
  "github.com/Masterminds/sprig"
  "html/template"
  "bytes"
)

func defineDatabaseRoutes(app *iris.Application)  {
  app.Get("/sql/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    fmt.Printf("Username: [%s]\n", s.GetString("username"))
    fmt.Printf("Getting full path %s\n", fullpath)
    fmt.Printf("Getting relative path [%s]\n", relativepath)
    reader, err :=  abucket.ReadStream(relativepath)
    templatefilename := c.URLParam("template")
    sqlstring := ""
    if err != nil {
      fmt.Println(err)
    } else {
      b, err := ioutil.ReadAll(reader)
      if err != nil {
        fmt.Print(err)
      } else {
        sqlstring = string(b)
      }
    }
    dbinfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", "odoo", "odoo", "192.168.10.69", "fimim")
    //ds := "username:password@tcp(host:3306)/db"
    //db, err := sql.Open("mysql", ds)
    db, err := sql.Open("postgres", dbinfo)

    if err != nil {
      fmt.Println("sql.Open:", err)
    }

    theCase := "lower" // "lower", "upper", "camel" or the orignal case if this is anything other than these three

    // headers []string, data [][]string, error
    //headers, data, _ := gosqljson.QueryDbToArray(db, theCase, "SELECT ID,NAME FROM t LIMIT ?,?", 0, 3)
    //fmt.Println(headers)
    // ["id","name"]
    //fmt.Println(data)
    // [["0","Alicia"],["1","Brian"],["2","Chloe"]]

    // data []map[string]string, error
    fmt.Printf("SQLSTRING\n")
    fmt.Println(sqlstring)
    data, err := gosqljson.QueryDbToMap(db, theCase, sqlstring)
    if err != nil {
      fmt.Println("sql.Query:", err)
    }
    //fmt.Printf("%+v\n", data)
    b, err :=json.Marshal(data)
    if err != nil {
      fmt.Println("json.Marshal:", err)
    }
    if templatefilename == "" {
      c.Write(b)
    } else {
      tstr := ""
      templatebb, err :=abucket.ReadStream(templatefilename)
      if err == nil {
        templateb, err := ioutil.ReadAll(templatebb)
        if err == nil {
          tstr = string(templateb)
        }
      } else {
        fmt.Printf("Error reading template %s\n", err)
      }

      t, err := template.New("at").Funcs(sprig.FuncMap()).Parse(tstr)
      if err != nil {
        fmt.Printf("template error %s\n", err)
      }

      var tpl bytes.Buffer
      err = t.Execute(&tpl, data)
      if err != nil {
        fmt.Printf("template error %s\n", err)
      }
      c.Write(tpl.Bytes())
    }
    // [{"id":"0","name":"Alicia"},{"id":"1","name":"Brian"},{"id":"2","name":"Chloe"}]
   c.Next()
  })
}
