package main

import (
  "fmt"
  "github.com/kataras/iris"
  "io/ioutil"
  "path/filepath"
  "bytes"
  "encoding/json"
  "html/template"
  ttemplate "text/template"
  "github.com/Masterminds/sprig"
  "github.com/satori/go.uuid"
  "os"
  "os/exec"
  "path"
  //"strings"
)


func defineTemplateRoutes(app *iris.Application, cacheDir string)  {
  app.Get("/template/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    templatefilename := c.URLParam("template")

    tstr := ""
    m := map[string]interface{}{}
    fmt.Printf("Executing template %s\n", templatefilename)
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
    //t, err := template.New("at").Parse(tstr)
    if err != nil {
      fmt.Printf("template error %s\n", err)
    }

    bb, err :=abucket.ReadStream(relativepath)
    if err == nil {
      b, err := ioutil.ReadAll(bb)
      if err == nil {
        m["contents"] = string(b)
        err2 := json.Unmarshal(b, &m)
        if err2 != nil {
          fmt.Printf("Error unmarshaling json %s\n", err2)
        }
      }
    } else {
      fmt.Printf("Error reading path %s %s\n", relativepath , err)
    }
    m["fullpath"] = fullpath
    m["relativepath"] = relativepath
    base := filepath.Base(relativepath)
    m["base"] = base
    ext := filepath.Ext(relativepath)
    m["ext"] = ext
    m["basenoext"] = base[0:len(base)-len(ext)]
    var tpl bytes.Buffer
    err = t.Execute(&tpl, m)
    if err != nil {
      fmt.Printf("template error %s\n", err)
    }
    c.Write(tpl.Bytes())
    c.Next()
  })


  app.Post("/formexec/{path:path}", func (c iris.Context) {
    fmt.Println("HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH<<<<<<<<<<***************<<<<<<<<<<<*************<<<<<<<<<<<<*********")
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")

    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    if c.Request().Body == nil {
      fmt.Println("Request body nil")
      //return errors.New("unmarshal: empty body")
    }

    rawData, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
      fmt.Println(err)
      //return err
    }
    fmt.Println(string(rawData))
    m := map[string]interface{}{}
    err2 := json.Unmarshal(rawData, &m)
    if err2 != nil {
      fmt.Printf("Error unmarshaling json %s\n", err2)
    }
    m["fullpath"] = fullpath
    m["relativepath"] = relativepath
    base := filepath.Base(relativepath)
    m["base"] = base
    ext := filepath.Ext(relativepath)
    m["ext"] = ext
    m["basenoext"] = base[0:len(base)-len(ext)]
    m["username"]=s.GetString("username")
    m["remoteip"]=c.Request().Header.Get("X-Forwarded-For")
    m["realip"]=c.Request().Header.Get("X-Real-Ip")

    tstr := ""
    templatefilename := (m["data"].(map[string]interface{})["template"]).(string)
    fmt.Printf("Executing template %s\n", templatefilename)
    templatebb, err :=abucket.ReadStream(templatefilename)
    if err == nil {
      templateb, err := ioutil.ReadAll(templatebb)
      if err == nil {
        tstr = string(templateb)
      }
    } else {
      fmt.Printf("Error reading template %s\n", err)
    }

    t, err := ttemplate.New("at").Funcs(ttemplate.FuncMap(sprig.FuncMap())).Parse(tstr)
    //t, err := ttemplate.New("at").Parse(tstr)
    if err != nil {
      fmt.Printf("template error %s\n", err)
    }


    var tpl bytes.Buffer
    err = t.Execute(&tpl, m)
    if err != nil {
      fmt.Printf("template error %s\n", err)
    }
    fmt.Println(string(tpl.Bytes()))
    u1 := uuid.NewV4()
    myTargetDir := path.Join(path.Dir(relativepath), "_config", "logs")
    /*
    err = os.MkdirAll(myTargetDir, 0755)
    if err != nil {
      fmt.Printf("Error creating dir %s\n", err)
    }
    myTargetFormFile := path.Join(myTargetDir, base)+"."+u1.String()+".script"
    err = abucket.WriteStream(myTargetFormFile, strings.NewReader(string(tpl.Bytes())), 0644)
    if err != nil {
      fmt.Println(err)
    }
    */
    myLocalDir := path.Join(cacheDir, myTargetDir)
    fmt.Printf("Creating dir %s\n", myLocalDir)
    err=os.MkdirAll(myLocalDir, 0755)
    if err != nil {
      fmt.Println(err)
    }
    myLocalExeFile := path.Join(myLocalDir, base)+"."+u1.String()+".script"
    err = ioutil.WriteFile(myLocalExeFile, tpl.Bytes(), 0755)
    if err != nil {
      fmt.Println(err)
    }

    cmd := exec.Command("bash", "-c", myLocalExeFile)
    stdoutStderr, err := cmd.CombinedOutput()
    if err != nil {
      fmt.Println(err)
    }
    fmt.Printf("%s\n", stdoutStderr)

    c.Write(tpl.Bytes())

    /*
    redirect := "" //`/result?text=El formulari s'ha enviat correctament.`
    redirect2, _ := jsonparser.GetString(rawData, "data", "redirect")
    if redirect2 != "" {
      redirect = redirect2
    }
    wfis, err := getWorkflow(abucket, myTargetFormFile, "write", s.GetString("username"))
    if err != nil {
      fmt.Printf(">>>>>>>>>>>>>>>>>>> Error getting workflow %s\n", myTargetFormFile)
    } else {
      for _, wfi := range *wfis {
        wfi.Run("write", abucket)
      }
    }

    */

    c.Next()
  })
}
