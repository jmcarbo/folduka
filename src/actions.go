package main

import (
  ttemplate "html/template"
  "os"
  "os/exec"
  "net/http"
  "folduka/bucket"
  "github.com/jordan-wright/email"
  "github.com/mholt/archiver"
  "net/smtp"
  "bytes"
  "encoding/json"
  "path"
  "fmt"
  "io/ioutil"
  "strings"
  "path/filepath"
  "text/template"
  "github.com/kataras/iris"
  "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
  "github.com/pdfcpu/pdfcpu/pkg/api"
  "strconv"
  "github.com/Masterminds/sprig"
  "github.com/rylio/ytdl"
)

type Parameter struct {
  Name string
  Label string
  Default string
}

type Action struct {
  Id string
  Type string
  Mask string
  Icon string
  Label string
  Parameters []Parameter
  Action string
  Prompt string
  Template string
}

func (a *Action) GetParameterDefault(parameter string) string {
  for _, p := range a.Parameters {
    if p.Name == parameter {
      return p.Default
    }
  }
  return ""
}

func (a *Action) RunWorkflowInstance(wfi *WorkflowInstance, abucket bucket.Bucket) {
  fmt.Printf("Running workflow %s with action %s\n", wfi.UUIDStr, a.Id)
  switch a.Action {
  case "email":
    originalfilename := wfi.CurrentPath
    m := map[string]interface{}{}
    bb, err :=abucket.ReadStream(originalfilename)
    if err == nil {
      b, err := ioutil.ReadAll(bb)
      if err == nil {
        m["contents"] = string(b)
        err2 := json.Unmarshal(b, &m)
        if err2 != nil {
          fmt.Printf("Error unmarshaling json %s\n", err2)
        }
      }
    }

    // Expand default parameters
    for _, p := range a.Parameters {
      m[p.Name] = p.Default
    }

    if originalfilename != "" {
      m["fullpath"] = originalfilename
      m["relativepath"] = originalfilename
    } else {
      m["fullpath"] = originalfilename
      m["relativepath"] = originalfilename
    }
    base := filepath.Base(originalfilename)
    m["base"] = base
    ext := filepath.Ext(originalfilename)
    m["ext"] = ext
    m["basenoext"] = base[0:len(base)-len(ext)]
    m["originalfilename"] = originalfilename
    mymail := email.NewEmail()
    mymail.From = "\"6imim\" <6imim@imim.es>"
    totemplate := m["to"].(string)
    ttotemplate, err := template.New("totemplate").Parse(totemplate)
    var totpl bytes.Buffer
    err = ttotemplate.Execute(&totpl, m)
    mymail.To = strings.Split(totpl.String(), ",")

    subjecttemplate := m["subject"].(string)
    tsubjecttemplate, err := template.New("subjecttemplate").Parse(subjecttemplate)
    var subjecttpl bytes.Buffer
    err = tsubjecttemplate.Execute(&subjecttpl, m)
    mymail.Subject = subjecttpl.String()
    //e.Text = []byte(fmt.Sprintf("You have asked for a password reset IMIM. Please enter reset code: %s when requested.", t.ResetCode))

    emailtemplate := m["email"].(string)
    temailtemplate, err := template.New("emailtemplate").Parse(emailtemplate)
    var emailtpl bytes.Buffer
    err = temailtemplate.Execute(&emailtpl, m)
    mymail.HTML = []byte(emailtpl.String())

    //err:=m.Send("hermes4.imim.es:587", smtp.PlainAuth("", "jmcarbo", "hvfcID24r", "hermes4.imim.es"))
    //auth := LoginAuth("6imim@imim.es", "j(6}@zCh)MA]utYLnJ1Rb")
    err=mymail.Send("mailuka.imim.science:587",
    smtp.PlainAuth("", "6imim", "j(6}@zCh)MA]utYLnJ1Rb", "mailuka.imim.science"))
    if err!=nil {
      fmt.Printf("Email error %s\n", err)
    } else {
      fmt.Println("Message sent correctly")
    }

  }
}

func getActions(abucket bucket.Bucket, mypath string) ([]Action, error) {
  fmt.Printf("Getting actions from bucket %s with path %s\n", abucket.Prefix(), mypath)
  actionsfilename := GetConfigFilename(abucket, mypath, "actions.json")
  r, err := abucket.ReadStream(actionsfilename)
  if err != nil {
    return nil, err
  }

  buf := new(bytes.Buffer)
  buf.ReadFrom(r)
  actions := []Action{}
  err = json.Unmarshal(buf.Bytes(), &actions)
  if err != nil {
    return nil, err
  }
  return actions, nil
}

func defineGetActions(app *iris.Application, cacheDir string) {
  app.Get("/action/{path:path}", func (c iris.Context) {
    fmt.Printf("**************++++++++++++++++++++******************+++++++++++++++++************\n")
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    action := c.URLParam("action")
    actionid := c.URLParam("actionid")
    actions, err := getActions(abucket, relativepath)
    if err != nil {
      fmt.Printf("Error getting actions %s\n", err)
    }
    fmt.Printf("Actions %+v\n", actions)
    fmt.Printf("Searching for action %s id: %s %+v\n", action, actionid, c.URLParams())
    switch action {
    case "archive":
    case "unarchive":
      myfilepath := path.Join(cacheDir, fullpath)
      os.MkdirAll(filepath.Dir(myfilepath), 0755)
      os.MkdirAll(path.Join(filepath.Dir(myfilepath), "outDir"), 0755)
      fmt.Printf("======%s\n=========%s\n", myfilepath, fullpath)
      abucket.Download(path.Join(relativepath), myfilepath)
      err := archiver.Unarchive(myfilepath, "outDir")
      if err != nil {
        fmt.Printf("unarchive error %s\n", err)
      }
    case "link":

      m := map[string]interface{}{}
      /*
      if err := json.Unmarshal([]byte(jsondata), &m); err != nil {
        panic(err)
      } 
      */
      myAction := Action{}
      fmt.Printf("Searching for actions id: %s\n", actionid)
      for _, a := range actions {
        if a.Id == actionid {

          myAction = a
          fmt.Printf(">>>>>>>>>>>>>> Found myAction %+v\n", myAction)
          for _, p := range a.Parameters {
            m[p.Name] = c.URLParam(p.Name)
          }
        }
      }

      tstr := myAction.Template
      t, err := ttemplate.New("at").Parse(tstr)
      if err != nil {
        fmt.Printf("template error %s\n", err)
      }

      bb, err :=abucket.ReadStream(relativepath)
      if err == nil {
        b, err := ioutil.ReadAll(bb)
        if err == nil {
          m["contents"] = string(b)
        }
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
      fmt.Printf("Redirecting to %s\n", tpl.String())
      c.Redirect(tpl.String())
      c.Next()
      return
    }
    c.Redirect(path.Join("/", fullpath))
    c.Next()
  })
}

func definePostActions(app *iris.Application, cacheDir string) {
  app.Post("/action/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    action := c.FormValue("action")
    actionid := c.FormValue("actionid")
    actions, _ := getActions(abucket, relativepath)
    fmt.Printf("Posting to %s\n", fullpath)
    switch action {
    case "mkdir":
      param1 := c.FormValue("folderName")
      fmt.Printf("Creating path %s\n", param1)
      abucket.Make(path.Join(relativepath, param1))
    case "copy":
      originalfilename := c.FormValue("originalfilename")
      newfilename := c.FormValue("newfilename")
      fmt.Printf("------- >Copying: %s %s %s\n", originalfilename, newfilename, relativepath)
      var err error
      if relativepath == "" {
        err=abucket.Copy(originalfilename, newfilename, false)
      }else {
        err =abucket.Copy(path.Join(relativepath, originalfilename), path.Join(relativepath, newfilename), false)
      }
      if err != nil {
        fmt.Printf("------- >ERROR Renaming %s %s %s %s\n", originalfilename, newfilename, relativepath, err)
      }
    case "rename":
      originalfilename := c.FormValue("originalfilename")
      newfilename := c.FormValue("newfilename")
      fmt.Printf("------- >Renaming %s %s %s\n", originalfilename, newfilename, relativepath)
      var err error
      if relativepath == "" {
        err =abucket.Rename(originalfilename, newfilename, false)
      }else {
        err = abucket.Rename(path.Join(relativepath, originalfilename), path.Join(relativepath, newfilename), false)
      }
      if err != nil {
        fmt.Printf("------- >ERROR Renaming %s %s %s %s\n", originalfilename, newfilename, relativepath, err)
      }
    case "ocrpdf":
      go func() {
        originalfilename := c.FormValue("originalfilename")
        myfilepath := path.Join(cacheDir, fullpath, originalfilename)
        os.MkdirAll(filepath.Dir(myfilepath), 0755)
        os.MkdirAll(path.Join(filepath.Dir(myfilepath), "outDir"), 0755)
        fmt.Printf("======%s\n=========%s\n", myfilepath, fullpath)
        abucket.Download(path.Join(relativepath, originalfilename), myfilepath)
        PutQueue(myfilepath+";"+originalfilename+";"+abucket.Name())
      }()
    case "splitpdf":
      originalfilename := c.FormValue("originalfilename")
      splitcountstr := c.FormValue("splitcount")
      splitcount, _ := strconv.Atoi(splitcountstr)
      config := pdfcpu.NewDefaultConfiguration()
      myfilepath := path.Join(cacheDir, fullpath, originalfilename)
      os.MkdirAll(filepath.Dir(myfilepath), 0755)
      os.MkdirAll(path.Join(filepath.Dir(myfilepath), "outDir"), 0755)
      fmt.Printf("======%s\n=========%s\n", myfilepath, fullpath)
      abucket.Download(path.Join(relativepath, originalfilename), myfilepath)
      err := api.SplitFile(myfilepath, path.Join(filepath.Dir(myfilepath),"outDir"), splitcount, config)
      if err != nil {
        fmt.Println(err)
      }

      files, err := ioutil.ReadDir(path.Join(filepath.Dir(myfilepath),"outDir"))
      if err != nil {
        fmt.Println(err)
      }
      for _, f := range files {
        abucket.Upload(path.Join(filepath.Dir(myfilepath), "outDir", f.Name()), path.Join(relativepath, f.Name()))

      }

    case "decorate":
      targetform := c.FormValue("targetForm")
      originalfilename := c.FormValue("originalfilename")
      c.Redirect(fmt.Sprintf("/formfill/%s.json?form=%s&decorate=%s", path.Join(fullpath, originalfilename), targetform, path.Join("/", fullpath, originalfilename)))
      c.Next()
      return
    case "formfill":
      targetform := c.FormValue("targetForm")
      originalfilename := c.FormValue("originalfilename")
      roottag := c.FormValue("roottag")
      if roottag != "" {
        c.Redirect(fmt.Sprintf("/formfill/%s?form=%s&root-tag=%s", path.Join(fullpath, originalfilename), targetform, roottag))
      } else {
        c.Redirect(fmt.Sprintf("/formfill/%s?form=%s", path.Join(fullpath, originalfilename), targetform))
      }
      c.Next()
      return
    case "post":
      originalfilename := c.FormValue("originalfilename")
      m := map[string]interface{}{}

      bb, err :=abucket.ReadStream(path.Join(relativepath, originalfilename))
      if err == nil {
        b, err := ioutil.ReadAll(bb)
        if err == nil {
          m["contents"] = string(b)
          err2 := json.Unmarshal(b, &m)
          if err2 != nil {
            fmt.Printf("POST Error unmarshaling json %s\n", err2)
          }
        }
      }

      if originalfilename != "" {
        afull:=path.Join(fullpath, originalfilename)
        arel :=path.Join(relativepath, originalfilename)

        m["fullpath"] = afull
        m["relativepath"] = arel
        base := filepath.Base(arel)
        m["base"] = base
        ext := filepath.Ext(arel)
        m["ext"] = ext
        m["basenoext"] = base[0:len(base)-len(ext)]
      } else {
        m["fullpath"] = fullpath
        m["relativepath"] = relativepath
        base := filepath.Base(relativepath)
        m["base"] = base
        ext := filepath.Ext(relativepath)
        m["ext"] = ext
        m["basenoext"] = base[0:len(base)-len(ext)]
      }
      m["originalfilename"] = originalfilename
      fmt.Printf("%+v\n", m)

      myAction := Action{}
      for _, a := range actions {
        if a.Id == actionid {
          myAction = a
          for _, p := range a.Parameters {
            m[p.Name] = c.FormValue(p.Name)
          }
        }
      }
      fmt.Printf("POST action %+v\n", myAction)

      totemplate := myAction.Template
      targeturl := c.FormValue("TargetURL")
      format := c.FormValue("Format")
      if format == "" {
        format = "application/json"
      }
      ttotemplate, err := template.New("totemplate").Parse(totemplate)
      var totpl bytes.Buffer
      err = ttotemplate.Execute(&totpl, m)
      if err != nil {
        fmt.Printf("POST error parsing template %s\n", err)
      }
      fmt.Printf("POST %s %s %s\n", targeturl, format, totpl.String())
      http.Post(targeturl, format, bytes.NewReader(totpl.Bytes()))
    case "email":
      originalfilename := c.FormValue("originalfilename")
      m := map[string]interface{}{}

      bb, err :=abucket.ReadStream(path.Join(relativepath, originalfilename))
      if err == nil {
        b, err := ioutil.ReadAll(bb)
        if err == nil {
          m["contents"] = string(b)
          err2 := json.Unmarshal(b, &m)
          if err2 != nil {
            fmt.Printf("Error unmarshaling json %s\n", err2)
          }
        }
      }

      if originalfilename != "" {
        m["fullpath"] = path.Join(fullpath, originalfilename)
        m["relativepath"] = path.Join(relativepath, originalfilename)
      } else {
        m["fullpath"] = fullpath
        m["relativepath"] = relativepath
      }
      base := filepath.Base(relativepath)
      m["base"] = base
      ext := filepath.Ext(relativepath)
      m["ext"] = ext
      m["basenoext"] = base[0:len(base)-len(ext)]
      m["originalfilename"] = originalfilename
      mymail := email.NewEmail()
      mymail.From = "\"6imim\" <6imim@imim.es>"
      totemplate := c.FormValue("to")
      ttotemplate, err := template.New("totemplate").Parse(totemplate)
      var totpl bytes.Buffer
      err = ttotemplate.Execute(&totpl, m)
      mymail.To = strings.Split(totpl.String(), ",")

      subjecttemplate := c.FormValue("subject")
      tsubjecttemplate, err := template.New("subjecttemplate").Parse(subjecttemplate)
      var subjecttpl bytes.Buffer
      err = tsubjecttemplate.Execute(&subjecttpl, m)
      mymail.Subject = subjecttpl.String()
      //e.Text = []byte(fmt.Sprintf("You have asked for a password reset IMIM. Please enter reset code: %s when requested.", t.ResetCode))

      emailtemplate := c.FormValue("email")
      temailtemplate, err := template.New("emailtemplate").Parse(emailtemplate)
      var emailtpl bytes.Buffer
      err = temailtemplate.Execute(&emailtpl, m)
      mymail.HTML = []byte(emailtpl.String())

      //err:=m.Send("hermes4.imim.es:587", smtp.PlainAuth("", "jmcarbo", "hvfcID24r", "hermes4.imim.es"))
      //auth := LoginAuth("6imim@imim.es", "j(6}@zCh)MA]utYLnJ1Rb")
      err=mymail.Send("mailuka.imim.science:587",
      smtp.PlainAuth("", "6imim", "j(6}@zCh)MA]utYLnJ1Rb", "mailuka.imim.science"))
      if err!=nil {
        fmt.Printf("Email error %s\n", err)
      } else {
        fmt.Println("Message sent correctly")
      }
    case "link":
      originalfilename := c.FormValue("originalfilename")

      m := map[string]interface{}{}
      /*
      if err := json.Unmarshal([]byte(jsondata), &m); err != nil {
        panic(err)
      } 
      */

      bb, err :=abucket.ReadStream(relativepath)
      if err == nil {
        b, err := ioutil.ReadAll(bb)
        if err == nil {
          m["contents"] = string(b)
        }
      }

      myAction := Action{}
      for _, a := range actions {
        if a.Id == actionid {
          myAction = a
          for _, p := range a.Parameters {
            m[p.Name] = c.FormValue(p.Name)
          }
        }
      }

      tstr := myAction.Template
      t, err := ttemplate.New("").Parse(tstr)
      if err != nil {
        fmt.Printf("template error %s\n", err)
      }

      m["originalfilename"] = originalfilename
      if originalfilename != "" {
        m["fullpath"] = path.Join(fullpath, originalfilename)
        m["relativepath"] = path.Join(relativepath, originalfilename)
      } else {
        m["fullpath"] = fullpath
        m["relativepath"] = relativepath
      }
      base := filepath.Base(relativepath)
      m["base"] = base
      ext := filepath.Ext(relativepath)
      m["ext"] = ext
      m["basenoext"] = base[0:len(base)-len(ext)]
      var tpl bytes.Buffer
      err = t.Execute(&tpl, m)
      c.Redirect(tpl.String())
      c.Next()
      return
    case "youtube":
      tag := c.FormValue("tag")
      vid, err := ytdl.GetVideoInfoFromID(tag)
      if err != nil {
        fmt.Println("Failed to get video info")
      } else {
        go func() {
          vidfile:=vid.Title + ".mp4"
          tempdir, _ := ioutil.TempDir(cacheDir, "VIDEO")
          tempscript, _ := ioutil.TempFile(tempdir, "VIDEO")
          fname := tempscript.Name()
          vid.Download(vid.Formats[0], tempscript)
          tempscript.Close()
          //ffmpeg -i Learn\ Bash\ Scripts\ -\ Tutorial.mp4 -pix_fmt yuv420p output.mp4
          cmd := exec.Command("ffmpeg", "-i", fname, "-pix_fmt", "yuv420p", "-f", "mp4", "-vcodec", "libx264", "-preset", "fast", "-profile:v", "main", "-acodec", "aac" , path.Join(tempdir,vidfile))
          msg, err := cmd.CombinedOutput()
          if err != nil {
            fmt.Printf("Error converting video %s %s\n", err, msg)
          } else {
            abucket.Upload(path.Join(tempdir,vidfile), path.Join(relativepath, vidfile))
          }
        }()
     }
    case "runscript":
      fmt.Printf("Calling runscript\n")
      originalfilename := c.FormValue("originalfilename")

      m := map[string]interface{}{}

      bb, err :=abucket.ReadStream(relativepath)
      if err == nil {
        b, err := ioutil.ReadAll(bb)
        if err == nil {
          m["contents"] = string(b)
        }
        if err := json.Unmarshal(b, &m); err != nil {
        // ignore if no json data
        }
      }


      myAction := Action{}
      for _, a := range actions {
        if a.Id == actionid {
          myAction = a
          for _, p := range a.Parameters {
            m[p.Name] = c.FormValue(p.Name)
          }
        }
      }

      m["originalfilename"] = originalfilename
      if originalfilename != "" {
        m["fullpath"] = path.Join(fullpath, originalfilename)
        m["relativepath"] = path.Join(relativepath, originalfilename)
      } else {
        m["fullpath"] = fullpath
        m["relativepath"] = relativepath
      }
      base := filepath.Base(relativepath)
      m["base"] = base
      ext := filepath.Ext(relativepath)
      m["ext"] = ext
      m["basenoext"] = base[0:len(base)-len(ext)]

      var tpl bytes.Buffer
      tstr := myAction.Template
      t, err := ttemplate.New("").Funcs(sprig.FuncMap()).Parse(tstr)
      if err != nil {
        fmt.Printf("template error %s\n", err)
      } else {
        err = t.Execute(&tpl, m)
        if err == nil {
          language := c.FormValue("language")
          switch language {
          case "shell":
            tempscript, _ := ioutil.TempFile(cacheDir, "SCRIPT")
            fname := tempscript.Name()
            tempscript.Close()
            ioutil.WriteFile(fname, tpl.Bytes(), 0700)
            os.Chmod(fname, 0777)
            cmd := exec.Command("bash", "-c", fname)
            msg, err := cmd.CombinedOutput()
            if err != nil {
              fmt.Printf("Error executing script %s %s\n", fname, msg)
            }
          }
        }
      }
    }
    fmt.Printf("REdirecting to %s\n", fullpath)
    c.Redirect(path.Join("/", fullpath))
    c.Next()
  })
}
