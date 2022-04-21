package main


import (
  "sync"
  "regexp"
  "github.com/bitly/go-simplejson"
  "os/exec"
  //jqpipe "github.com/threatgrid/jqpipe-go"
  //"github.com/cosmos72/gomacro/fast"
  //"reflect"
  //"github.com/d5/tengo/script"
  "github.com/buger/jsonparser"
  "github.com/satori/go.uuid"
  "github.com/dchest/captcha"
  "github.com/lalamove/konfig"
  "github.com/lalamove/konfig/parser/kpjson"
  "github.com/lalamove/konfig/loader/klfile"
  "github.com/kataras/iris/v12"
  "github.com/kataras/iris/v12/sessions"
  "github.com/gorilla/securecookie"
  "github.com/iris-contrib/middleware/csrf"
  "github.com/casbin/casbin"
  "github.com/casbin/json-adapter"
  "gopkg.in/ldap.v2"
  "fmt"
  "log"
  "folduka/bucket"
  "html/template"
  "io/ioutil"
  //"io"
  "os"
  "path"
  "path/filepath"
  "strings"
  "bytes"
  "folduka/user"
  "sort"
  //"time"
  "encoding/json"
  ttemplate "text/template"
  "github.com/bregydoc/gtranslate"
  //"github.com/avct/uasurfer"
  "mime"
)


///Applications/LibreOffice.app/Contents/MacOS/soffice --headles --convert-to pdf *.*
var configFiles = []klfile.File{
  {
    Path: "./config/config.json",
    Parser: kpjson.Parser,
  },
}

func init() {
  konfig.Init(konfig.DefaultConfig())
}

func before(c iris.Context) {
    // fmt.Printf("Connecting from >>>>>>>>>>>>>>>>> %s\n", c.RemoteAddr())
    // [...]
    // fmt.Println("Before")
  fmt.Printf("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^\n")
  fmt.Printf("%+v\n", c.URLParams())
    s := mySessions.Start(c)

    // path constraint for user
    rootPath := s.GetString("RootPath")

    currentPath := c.Path()

    // check for public path
    isPublic := false
    public := []string{ "/css/", "/js/", "/img/", "/login", "/captcha", "/svg/", "/id", "/uuid", "/ViewerJS/", "/pdf/", "/build", "/favicon.ico" }
    for _, k := range public {
      if strings.HasPrefix(currentPath, k) {
        isPublic = true
        break
      }
    }

    if isPublic {
      c.Next()
      return
    }
    // end of public check

    // onlyofficesave
    if strings.HasPrefix(currentPath, "/onlyofficesave") {
      // TODO: allow post without CSRF token
      csrf.UnsafeSkipCheck(c)
    }

    //fmt.Printf("%+v\n", c.GetReferrer())
    //fmt.Printf("%+v\n", c.Path()) 
    //fmt.Printf("%+v\n", c.RemoteAddr())
    n:=s.GetString("name")
    fullpath := c.Params().Get("path")

    // Get bucket
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)

    //fmt.Printf("%+v %s\n", abucket, relativepath)
    //fmt.Printf("Bucket %s\n", abucket.Name())
    curfilestat, err := abucket.Stat(relativepath)


    dir := "/"
    // if no err in filepath
    if err == nil {
      if curfilestat.IsDir() {
        dir = relativepath
      } else {
        dir = filepath.Dir(relativepath)
      }
    }

    // Get policyfile
    policyfilename := GetPolicyFilename(abucket, dir)
    //fmt.Printf("Policyfilename [%s]\n", policyfilename)
    //policyfilename := path.Join(dir, "_config", "policy.json")
    policyfile, err := abucket.Stat(policyfilename)
    if err == nil && policyfile != nil {
      fmt.Printf(">>>>>>>>>>>>>>>> %s\n", policyfile.Name())
    }

    /* REDUNDANT ?
    afile2, err := abucket.Stat(relativepath)
    if err == nil && afile2 != nil {
      fmt.Printf(">>>>>>>>>>>>>>>> %s\n", afile2.Name())
    } else {
      fmt.Printf("File not found !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
    }
    */

    // see if local enforcer
    var e *casbin.Enforcer
    hasEnforcer := false
    if policyfile != nil {
      bb, err :=abucket.ReadStream(policyfilename)
      if err != nil {
        fmt.Println(err)
      } else {
        casbinpolicybytes, err := ioutil.ReadAll(bb)
        if err != nil {
          fmt.Println(err)
        } else {
          a := jsonadapter.NewAdapter(&casbinpolicybytes)
          e = casbin.NewEnforcer("./config/model.conf", a)
          err = e.LoadPolicy()
          if err != nil {
            fmt.Println(err)
          }
          e.AddPermissionForUser("jmcarbo@imim.es", ".*", "read", "allow")
          permissions := e.GetPermissionsForUser("jmcarbo@imim.es")
          fmt.Printf("PERMISSIONS ########################### %+v\n", permissions)
          c.Values().Set("Enforcer", e)
          hasEnforcer = true
        }
      }
    }

    // load default enforcer
    if !hasEnforcer {
      casbinpolicybytes, err := ioutil.ReadFile("./config/policy.json")
      if err != nil {
        fmt.Println(err)
      } else {
        a := jsonadapter.NewAdapter(&casbinpolicybytes)
        e = casbin.NewEnforcer("./config/model.conf", a)
        err = e.LoadPolicy()
        if err != nil {
          fmt.Println(err)
        } else {
          c.Values().Set("Enforcer", e)
          hasEnforcer = true
        }
      }
    }

    username := s.GetString("username")
    fmt.Printf("=========== n: [%s] rootPath: [%s] isPublic: [%t] username: [%s] fullpath: [%s] currentPath [%s] relativepath: [%s]\n", 
      n, rootPath, isPublic, username, fullpath, currentPath, relativepath)

    if n == "iris" {
      if e != nil {
        e.AddRoleForUser(username, "authorized")
        if s.GetString("imimuser") == "true" {
          e.AddRoleForUser(username, "imimuser")
        }
      }
      fmt.Println("---> Registered user")
      // Registered user
      // Check if RootPath is compatible
      if rootPath != "" && !strings.HasPrefix(fullpath, rootPath) && !strings.HasPrefix(currentPath, "/logout") {
        fmt.Println("---> Registered user rootpath active")
        /*
        if fullpath == "" {
          c.Redirect("/login")
          c.Next()
          return
        } else {
          c.Redirect(fmt.Sprintf("/login?path=%s", c.Path()))
          c.Next()
          return
        }
        */
      }
    } else if n == "" && !isPublic {
      isAuthorized := false
      if e != nil {
        // check if anonymous is authorized
        isAuthorized = e.Enforce("anonymous", relativepath, "read")
        fmt.Printf("Relativapath: %s isauthorized: %t\n", relativepath, isAuthorized)
        if isAuthorized {
          s.Set("username", "anonymous")
          c.Values().Set("username", "anonymous")
        }

        // check if token authorized
        tokenstr:=c.URLParam("token")
        if tokenstr != "" {
          cfn := GetConfigFilename(abucket, relativepath, path.Join("permissions", filepath.Base(relativepath) + "." + tokenstr))
          if cfn != "" {
            isAuthorized = true
            c.Values().Set("token", tokenstr)
          }
        }
      }

      if !isAuthorized {
        if fullpath == "" {
          c.Redirect("/login")
          c.Next()
          return
        } else {
          c.Redirect(fmt.Sprintf("/login?path=%s", c.Path()))
          c.Next()
          return
        }
      }
    }


    c.Next()
}

func after(ctx iris.Context) {
           // [...]
       fmt.Println("After")
}

var mySessions *sessions.Sessions

var mybucket bucket.Bucket
var uploadsDir = "./uploads"
var bucketNames []string
var buckets = map[string]bucket.Bucket{}
var bucketSync = sync.Mutex{}

func init() {

}

func breadcrumbPath(currentPath, defaultBucket string) template.HTML {
  fmt.Printf("Breadcrumb [%s]\n", currentPath)
  if currentPath == "/" || currentPath == "" {
    return template.HTML("/")
  }

  abucket, relativepath := LoadBucketFromPath(currentPath, defaultBucket)
  info, err := abucket.Stat(relativepath)
  if err != nil {
    return template.HTML("")
  }
  elements := strings.Split(currentPath, "/")

  lastElement := len(elements) -1

  if !info.IsDir() {
    lastElement--
  }
  if lastElement == -1 {
    return template.HTML(elements[0])
  }
  elementBeforeLast := lastElement-1
  fmt.Printf("%d %d\n", elementBeforeLast, lastElement)
  breadcrumb := ""
  if elementBeforeLast == - 1 {
    breadcrumb = fmt.Sprintf(`<a href="/"><img src="/svg/arrow-thick-left.svg" width="16" alt="arrow-left"></a> %s`, elements[0])
  } else {
    breadcrumb = fmt.Sprintf(`<a href="/%s"><img src="/svg/arrow-thick-left.svg" width="16" alt="arrow-left"></a>&nbsp;`, strings.Join(elements[0:elementBeforeLast+1], "/"))
    for i:= 0 ; i<=elementBeforeLast; i++ {
      breadcrumb = breadcrumb + fmt.Sprintf(`<a href="/%s">%s</a> /`, strings.Join(elements[0:i+1],"/"), elements[i])
    }
    breadcrumb = breadcrumb + " " + elements[lastElement]
  }
  return template.HTML(breadcrumb)
}

func LoadBucket(bucketname string) bucket.Bucket {
  bucketSync.Lock()
  defer bucketSync.Unlock()
  return buckets[bucketname]
}

func InitializeBuckets() {
  for _, k := range bucketNames {
    bucketname := k
      switch konfig.String(k+".class") {
      case "local":
        buckets[k]= bucket.NewLocalBucket(konfig.String(bucketname+".name"),konfig.String(bucketname+".root"),
          konfig.String(bucketname+".host"), konfig.String(bucketname+".user"), konfig.String(bucketname+".password") )
        buckets[k].SetPrefix(k)
      case "webdav":
        buckets[k]= bucket.NewDAVBucket(konfig.String(bucketname+".name"),konfig.String(bucketname+".root"),
          konfig.String(bucketname+".host"), konfig.String(bucketname+".user"), konfig.String(bucketname+".password") )
        buckets[k].SetPrefix(k)
      case "imap":
        buckets[k]= bucket.NewImapBucket(konfig.String(bucketname+".name"),konfig.String(bucketname+".root"),
          konfig.String(bucketname+".host"), konfig.String(bucketname+".user"), konfig.String(bucketname+".password") )
        buckets[k].SetPrefix(k)
      }
  }
}

func getDefaultBucket(s *sessions.Session) string {
  defaultBucket := konfig.String("default_bucket")
  dbu:=s.GetString("default_bucket")
  if dbu != "" {
    defaultBucket = dbu
  }
  return defaultBucket
}

func LoadBucketFromPath(targetpath, defaultBucket string) (bucket.Bucket, string) {
  bucketSync.Lock()
  defer bucketSync.Unlock()
  for _, k := range bucketNames {
    j := k
    if !strings.HasPrefix(j, "/") {
      j = "/"+j
    }
    if strings.HasPrefix(targetpath, j) || strings.HasPrefix(targetpath, k) {
      fmt.Printf("Found prefix: [%s] [%s] [%s]\n", k, targetpath, strings.TrimPrefix(strings.TrimPrefix(targetpath, k), "/"))

      return buckets[k], strings.TrimPrefix(strings.TrimPrefix(targetpath, k), "/")
    } else {
      fmt.Printf("[%s] has no prefix %s\n", targetpath, k)
    }
  }
  return buckets[defaultBucket], targetpath
}


func main() {
  StartQueue()
  defer StopQueue()
  os.MkdirAll("./uploads", 0755)
  os.MkdirAll("./cache", 0755)
  cacheDir:="./cache"
  // load from json file 
  konfig.RegisterLoaderWatcher(
    klfile.New(&klfile.Config{
      Files:    configFiles,
      Watch:    true,
    }),
    // optionally you can pass config hooks to run when a file is changed
    func(c konfig.Store) error {
      bucketNames = konfig.StringSlice("buckets")
      InitializeBuckets()
      defaultBucket := konfig.String("default_bucket")
      if  defaultBucket == "" {
        mybucket = LoadBucket("bucket")
      } else {
        mybucket = LoadBucket(defaultBucket)
      }
      return nil
    },
  )

  if err := konfig.LoadWatch(); err != nil {
    log.Fatal(err)
  }
  bucketNames = konfig.StringSlice("buckets")
  InitializeBuckets()
  defaultBucket := konfig.String("default_bucket")
  if  defaultBucket == "" {
    mybucket = LoadBucket("bucket")
  } else {
    mybucket = LoadBucket(defaultBucket)
  }
  onlyofficeserver := konfig.String("onlyoffice")
  servername := konfig.String("servername")

  if onlyofficeserver == "" {
    onlyofficeserver = "onlyoffice.prod03.imim.science"
  }

  if servername == "" {
    servername = "folduka.imim.science"
  }


  casbinpolicybytes, err := ioutil.ReadFile("./config/policy.json")
  if err != nil {
    fmt.Println(err)
  }
  a := jsonadapter.NewAdapter(&casbinpolicybytes)
  e := casbin.NewEnforcer("./config/model.conf", a)
  err = e.LoadPolicy()
  if err != nil {
    fmt.Println(err)
  }
  fmt.Printf("%v\n", e.Enforce("jmcarbo@imim.es", "Gdata1", "read"))
  fmt.Printf("%v\n", e.Enforce("jmcarbo@imim.es", "data1", "read"))
  fmt.Printf("%v\n", e.Enforce("jmcarbo@imim.es", "data1", "write"))

  app := iris.Default()
  mime.AddExtensionType(".wasm","application/wasm")

  cookieName := "imimsciencefolduka"
  // AES only supports key sizes of 16, 24 or 32 bytes.
  // You either need to provide exactly that amount or you derive the key from what you type in.
  hashKey := []byte("9i5-456-323-secret-fash-key-here")
  blockKey := []byte("4o6-secret-58-characters-675-too")
  secureCookie := securecookie.New(hashKey, blockKey)

  protect := csrf.Protect([]byte("9AB0E421E53A477C084477AEA06096F5"),
  csrf.Secure(konfig.Bool("secure"))) // Defaults to true, but pass `false` while no https (devmode).

  mySessions = sessions.New(sessions.Config{
    Cookie:       cookieName,
    Encode:       secureCookie.Encode,
    Decode:       secureCookie.Decode,
    AllowReclaim: true,
  })
  app.Use(before)
  app.Use(protect)
  app.Done(after)
  app.StaticWeb("/js", "./public/js")
  app.StaticWeb("/css", "./public/css")
  app.StaticWeb("/img", "./public/img")
  app.StaticWeb("/fonts", "./public/fonts")
  app.StaticWeb("/svg", "./public/svg")
  app.StaticWeb("/pdf", "./public/pdf")
  app.StaticWeb("/build", "./public/build")
  app.StaticWeb("/ViewerJS", "./public/ViewerJS")
  app.StaticWeb("/acejs", "./public/acejs")
  // load templates
  tmpl:=iris.HTML("./templates", ".html")
  tmpl.AddFunc("regexMatch", func(regex string, s string) bool {
      match, _ := regexp.MatchString(regex, s)
      return match
  })

  tmpl.AddFunc("display", func(e bucket.Element) template.HTML {
    displayStr := ""
    //fmt.Printf("****************** %v\n", template.HTML(e.Display()))
    folder := e.Path()
    relativepath := path.Join(e.Path(), e.Name())
    display, err := getDisplay(e.Bucket(), folder)
    if err != nil {
      fmt.Printf("Error display %s\n", err)
      return template.HTML(displayStr)
    }
    //fmt.Printf("Display: %+v\n", display)
    for _, d := range display {
      m, _ := regexp.MatchString(d.Mask, relativepath)
      if m {
        //fmt.Printf("display template found %s\n", d.Mask)
        t, err := ttemplate.New("").Parse(d.Template)
        if err != nil {
          fmt.Printf("display template error %s\n", err)
          return template.HTML(displayStr)
        }

        m := map[string]interface{}{}

        bb, err :=e.Bucket().ReadStream(relativepath)
        if err == nil {
          b, err := ioutil.ReadAll(bb)
          if err == nil {
            m["contents"] = string(b)
            if err := json.Unmarshal([]byte(b), &m); err != nil {
              fmt.Printf("display parsing error %s in %s\n", err, e.Path())
            }
          }
        } else {
          fmt.Printf("display reading error %s in %s\n", err, e.Path())
        }
        m["fullpath"] = path.Join(e.Prefix(), e.Path(), e.Name()) 
        m["relativepath"] = relativepath
        base := filepath.Base(relativepath)
        m["base"] = base
        ext := filepath.Ext(relativepath)
        m["ext"] = ext
        m["basenoext"] = base[0:len(base)-len(ext)]
        //fmt.Printf(">>>>>>>> %+v\n", m)
        var tpl bytes.Buffer
        err = t.Execute(&tpl, m)
        if err != nil {
          fmt.Printf("display template execution error %s\n", err)
        }
        displayStr=tpl.String()
      }
    }
    return template.HTML(displayStr)
  })
  tmpl.AddFunc("isdir", func(e bucket.Element) bool {
    return e.IsDir()
  })
  tmpl.AddFunc("elementactions", func(enf *casbin.Enforcer, user string, e bucket.Element) template.HTML {
    mypath := ""
    if e.Prefix() == "" {
      mypath = path.Join(e.Path(), e.Name())
    } else {
      mypath = path.Join(e.Prefix(), e.Path(), e.Name())
    }
    actions := ""
    allow := enf.Enforce(user, mypath, "delete")
    if allow {
      delete_action := fmt.Sprintf(`<a href="/delete/%s"><img src="/svg/trash.svg" width="16" alt="trash"></a>`, mypath)
      actions = actions + delete_action
    }
    return template.HTML(actions)
  })
  tmpl.AddFunc("getaction", func(e bucket.Element) string {
    mypath := ""
    if e.Prefix() == "" {
      mypath = path.Join(e.Path(), e.Name())
    } else {
      mypath = path.Join(e.Prefix(), e.Path(), e.Name())
    }
    switch filepath.Ext(e.Name()) {
    case ".json":
      return fmt.Sprintf("/formfill/%s", mypath)
    case ".data":
      return fmt.Sprintf("/formfill/%s", mypath)
    case ".form":
      return fmt.Sprintf("/formfill/%s", mypath)
    default:
      return fmt.Sprintf("/download/%s", mypath)
    }
    return ""
  })
  tmpl.AddFunc("getextension", func(e bucket.Element) string {
    return filepath.Ext(e.Name())
  })
  tmpl.AddFunc("getname", func(e bucket.Element) string {
    return e.Name()
  })
  tmpl.AddFunc("getsize", func(e bucket.Element) string {
    return fmt.Sprintf("%-5d MB", e.Size()/1000000)
  })
  tmpl.AddFunc("getmodtime", func(e bucket.Element) string {
    return e.ModTime().Format("2006-01-02 15:04:05")
  })
  tmpl.AddFunc("getpath", func(e bucket.Element) string {
    a:=""
    if e.Prefix() == "" {
      a=path.Join("/", e.Path(), e.Name())
    }else {
      a=path.Join("/", e.Prefix(), e.Path(), e.Name())
    }
    fmt.Println(a)
    return a
  })
  tmpl.AddFunc("displayHeader", func(b bucket.Bucket) template.HTML {
    return template.HTML(b.DisplayHeader())
  })
  tmpl.AddFunc("displayFooter", func(b bucket.Bucket) template.HTML {
    return template.HTML(b.DisplayFooter())
  })
  tmpl.AddFunc("breadcrumbs", func(p, defaultBucket string) template.HTML {
    return breadcrumbPath(p, defaultBucket)
  })
  tmpl.AddFunc("jsSave", func(s string) template.JS {
    return template.JS(s)
  })
  tmpl.AddFunc("htmlSafe", func(s string) template.HTML {
    return template.HTML(s)
  })
  tmpl.AddFunc("enforce", func(e *casbin.Enforcer, user, resource, action string) bool {
    allow := e.Enforce(user, resource, action)
    fmt.Printf("Enforcer user: %s, resource %s action %s allow: %t\n", user, resource, action, allow)
    return allow
  })
  tmpl.AddFunc("getroles", func(e *casbin.Enforcer, user string) []string {
	  sa, _ := e.GetRolesForUser(user)
	  return sa
  })
  app.RegisterView(tmpl)

  //Digital signage
  stream(app, cacheDir)

  setupWebsocket(app)

  defineUtilFunctions(app)
  defineMailFunctions(app, cacheDir)
  defineTemplateRoutes(app, cacheDir)
  defineDatabaseRoutes(app)

  app.Get("/result", func(c iris.Context) {
    for k,v :=range c.URLParams() {
      c.ViewData(k, v)
    }
    c.View("result.html")
  })
  app.Get("/", func (c iris.Context) {
    c.ViewData("BucketName", mybucket.Name())
    c.ViewData("Bucket", mybucket)
    c.ViewData("Path", "/")
    c.ViewData("CSRFToken", csrf.Token(c))
    c.ViewData(csrf.TemplateTag, csrf.TemplateField(c))
    s := mySessions.Start(c)
    defaultBucket := getDefaultBucket(s)
    abucket, _ := LoadBucketFromPath("/", defaultBucket)
    c.ViewData("DefaultBucket", defaultBucket)
    c.ViewData("Username", s.GetString("username"))
    c.ViewData("BucketPrefix", abucket.Prefix())
    info, err := mybucket.Stat("/")
    if err != nil {
      fmt.Println(err)
    }
    c.ViewData("Info", info)
    if info.IsDir() {
      c.ViewData("IsDir", true)
        actions, err := getActions(abucket, "/")
        if err != nil {
          fmt.Printf("Error in actions: %s\n", err)
        }
        c.ViewData("actions", actions)
      els := abucket.List("/")
      sort.Sort(bucket.ElementByName(*els))
      c.ViewData("Elements", els)
    } else {
      c.ViewData("IsDir", false)
    }
    e2 := c.Values().Get("Enforcer")
    if e2 == nil {
      e2 = e
    }
    c.ViewData("Enforcer", e2)
    c.View("index.html")
    c.Next()
  })

  app.Get("/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    fmt.Printf("Username: [%s]\n", s.GetString("username"))
    fmt.Printf("Getting full path %s\n", fullpath)
    fmt.Printf("Getting relative path [%s]\n", relativepath)
    info, err := abucket.Stat(relativepath)
    if err != nil {
      fmt.Println(err)
      c.ViewData("IsDir", false)
    } else {
      c.ViewData("Info", info)
      if info != nil && info.IsDir() {
        // View folder
        c.ViewData("IsDir", true)
        actions, err := getActions(abucket, relativepath)
        if err != nil {
          fmt.Printf("Error in actions: %s\n", err)
        }
        c.ViewData("actions", actions)

        var els *[]bucket.Element
        search := c.URLParam("search")
        if search != "" {
          fmt.Printf("Searching %s %s\n", search, relativepath)
          els = abucket.Search(search, relativepath)
        } else {
          els = abucket.List(relativepath)
        }
        sort.Sort(bucket.ElementByName(*els))
        c.ViewData("Elements", els)
      } else {
        // View file
        //c.ViewData("IsDir", false)
        ext := filepath.Ext(fullpath)
        if ext != "" {
          ext = ext[1:]
        }

        token := c.Values().Get("token")

        if (strings.Contains("docx,xlsx,pptx", strings.ToLower(ext))) && (token == nil) {
          filename := filepath.Base(fullpath)
          c.ViewData("fullpath", path.Join(abucket.Prefix(),relativepath))
          c.ViewData("filename", filename)
          c.ViewData("filetype", ext)
          switch ext {
          case "pdf":
            c.ViewData("documentType", "text")
          case "docx":
            c.ViewData("documentType", "text")
          case "xlsx":
            c.ViewData("documentType", "spreadsheet")
          case "pptx":
            c.ViewData("documentType", "presentation")
          }
          u1 := uuid.NewV4()
          dir := filepath.Dir(relativepath)
          if dir == "" {
            dir = "/"
          }
          permdir := path.Join(dir, "_config", "permissions")
          err := abucket.Make(permdir)
          if err != nil {
            fmt.Printf(">>>>>> Error %s creating %s\n", err, permdir)
          }
          mypath:=path.Join(dir, "_config", "permissions", filename + "." + u1.String())
          err = abucket.WriteStream(mypath, strings.NewReader(string("")), 0644)
          c.ViewData("onlyofficetoken", u1.String())
          c.ViewData("onlyofficeserver", onlyofficeserver)
          c.ViewData("servername", servername)
          c.View("onlyofficeview.html")
        } else {
          myfilepath := path.Join(cacheDir, fullpath)
          os.MkdirAll(filepath.Dir(myfilepath), 0755)
          fmt.Printf("======%s\n=========%s\n", myfilepath, fullpath)
          abucket.Download(relativepath, myfilepath)
          c.Header("Access-Control-Allow-Origin", "*")
          fmt.Printf("/////////////////////////////////////////////////////////////////////////////\n")
          //myUA := c.GetHeader("User-Agent")
          // Parse() returns all attributes, including returning the full UA string last
          //ua := uasurfer.Parse(myUA)
          //fmt.Printf("%+v\n", ua)
          //if ua.Browser.Name.String() == "BrowserSafari" {
          //  c.ViewData("videofilename", fullpath)
          //  c.View("video.html")
          //} else {
	  if path.Ext(myfilepath) == ".wasm" {
	    c.ContentType("application/wasm")
	  }
          c.ServeFile(myfilepath, true)
          //}
        }
        c.Next()
        return
      }
    }
    c.ViewData("Username", s.GetString("username"))
    c.ViewData("DefaultBucket", defaultBucket)
    c.ViewData("BucketName", abucket.Name())
    c.ViewData("Bucket", abucket)
    c.ViewData("BucketPrefix", abucket.Prefix())
    c.ViewData("Path", fullpath)
    c.ViewData("RelativePath", relativepath)
    c.ViewData("CSRFToken", csrf.Token(c))
    c.ViewData(csrf.TemplateTag, csrf.TemplateField(c))
    e2 := c.Values().Get("Enforcer")
    if e2 == nil {
      fmt.Printf("----> Default enforcement\n")
      e2 = e
    }
    c.ViewData("Enforcer", e2)
    c.View("index.html")
    c.Next()
  })

  app.Get("/download/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    //ServeFile(filename string, gzipCompression bool) error
    info, err := abucket.Stat(relativepath)
    if err != nil {
      fmt.Println(err)
    } else {
      if info != nil && !info.IsDir() {
        myfilepath := path.Join(cacheDir, fullpath)
        os.MkdirAll(filepath.Dir(myfilepath), 0755)
        fmt.Printf("======%s\n=========%s\n", myfilepath, fullpath)
        abucket.Download(relativepath, myfilepath)
        c.SendFile(myfilepath, path.Base(relativepath))
      } else {
      }
    }
    c.Next()
  })

  app.Get("/sign/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    signature := c.URLParam("signature")
    password := c.URLParam("password")
    comments := c.URLParam("comments")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    //ServeFile(filename string, gzipCompression bool) error
    info, err := abucket.Stat(relativepath)
    if err != nil {
      fmt.Printf("Error signing %s %s\n", fullpath, err)
    } else {
      if info != nil && !info.IsDir() {
        myfilepath := path.Join(cacheDir, fullpath)
        mysignature := path.Join(cacheDir, signature)
        os.MkdirAll(filepath.Dir(myfilepath), 0755)
        os.MkdirAll(filepath.Dir(mysignature), 0755)
        fmt.Printf("Downloading pdf %s\n", relativepath)
        abucket.Download(relativepath, myfilepath)
        fmt.Printf("Downloading signature %s\n", signature)
        errsig := abucket.Download(signature, mysignature)
        if errsig != nil {
          fmt.Printf("Error downloading signature %s: %s\n", signature, errsig)
        }
        cmd := exec.Command("java", "-jar", "signpdf/PortableSigner.jar", "-n", "-s", mysignature, "-b", "es", "-c", comments, "-l", "location", "-t", myfilepath, "-o", myfilepath+".signed.pdf", "-p", password)
        msg, err := cmd.CombinedOutput()
        if err != nil {
          fmt.Printf("signing error ---> %s %s\n", err, msg)
        } else {
          fmt.Printf("Signed %s\n", msg)
          abucket.Upload(myfilepath + ".signed.pdf", relativepath + ".signed.pdf")
        }
      }
    }
    c.Redirect("/"+filepath.Dir(fullpath))
    c.Next()
  })

  app.Get("/display/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    //ServeFile(filename string, gzipCompression bool) error
    info, err := abucket.Stat(relativepath)
    if err != nil {
      fmt.Println(err)
    } else {
      if info != nil && !info.IsDir() {
        myfilepath := path.Join(cacheDir, fullpath)
        os.MkdirAll(filepath.Dir(myfilepath), 0755)
        fmt.Printf("======%s\n=========%s\n", myfilepath, fullpath)
        abucket.Download(relativepath, myfilepath)
        c.ServeFile(myfilepath, true)
      } else {
      }
    }
    c.Next()
  })
  app.Get("/edit/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    c.ViewData("BucketName", mybucket.Name())
    c.ViewData("Bucket", mybucket)
    c.ViewData("Path", c.Params().Get("path"))
    c.ViewData("Base", filepath.Base(fullpath))
    c.ViewData("CSRFToken", csrf.Token(c))
    ext := filepath.Ext(fullpath)
    c.ViewData("Ext", ext)
    bb, err :=abucket.ReadStream(relativepath)
    if err == nil {
      b, err := ioutil.ReadAll(bb)
      if err == nil {
        c.ViewData("content", string(b))
      }
    }


    if c.URLParam("editor") == "ace" {
      if c.URLParam("mode") != "" {
        c.ViewData("Mode", "ace/mode/" + c.URLParam("mode"))
      } else {
        c.ViewData("Mode", "ace/mode/text")
      }
      c.View("editace.html")
    }  else if (strings.Contains(".pdf,.doc,.docx,.xls,.xlsx,.pptx", ext)) {
      if ext != "" {
        ext = ext[1:]
      }
      filename := filepath.Base(fullpath)
      c.ViewData("fullpath", path.Join(abucket.Prefix(),relativepath))
      c.ViewData("filename", filename)
      c.ViewData("filetype", ext)
      switch ext {
      case "pdf":
        c.ViewData("documentType", "text")
      case "docx":
        c.ViewData("documentType", "text")
      case "xlsx":
        c.ViewData("documentType", "spreadsheet")
      case "pptx":
        c.ViewData("documentType", "presentation")
      }
      u1 := uuid.NewV4()
      dir := filepath.Dir(relativepath)
      if dir == "" {
        dir = "/"
      }
      permdir := path.Join(dir, "_config", "permissions")
      err := abucket.Make(permdir)
      if err != nil {
        fmt.Printf(">>>>>> Error %s creating %s\n", err, permdir)
      }
      mypath:=path.Join(dir, "_config", "permissions", filename + "." + u1.String())
      err = abucket.WriteStream(mypath, strings.NewReader(string("")), 0644)
      c.ViewData("onlyofficetoken", u1.String())
      c.ViewData("onlyofficeserver", onlyofficeserver)
      c.ViewData("servername", servername)
      c.View("onlyofficeedit.html")
    } else {
      switch ext {
      case ".md":
        c.View("editmd.html")
      case ".json":
        c.ViewData("Mode", "ace/mode/json")
        c.View("editace.html")
      case ".htm":
      case ".html":
        c.ViewData("Mode", "ace/mode/html")
        c.View("editck.html")
      case ".js":
        c.ViewData("Mode", "ace/mode/javascript")
        c.View("editace.html")
      default:
        if c.URLParam("mode") != "" {
          c.ViewData("Mode", "ace/mode/" + c.URLParam("mode"))
        } else {
          c.ViewData("Mode", "ace/mode/text")
        }
        c.View("editace.html")
      }
    }
    c.Next()
  })
  app.Post("/edit/{path:path}", func (c iris.Context) {
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
    err = abucket.WriteStream(relativepath, strings.NewReader(string(rawData)), 0644)
    if err != nil {
      fmt.Println(err)
    }
    c.JSON(map[string]string{"Result": "OK"})
    c.Next()
  })

  app.Post("/onlyofficesave/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    if c.Request().Body == nil {
      fmt.Println("_______________ Onlyoffice: Request body nil")
      //return errors.New("unmarshal: empty body")
    }

    rawData, err := ioutil.ReadAll(c.Request().Body)
    if err != nil {
      fmt.Println(fmt.Sprintf("___________________ Onlyoffice error reading body: %s\n", err))
      //return err
    }
    fmt.Println(string(rawData))
    myjson, err := simplejson.NewJson(rawData)
    if err != nil {
      fmt.Printf("Simplejson error in onlyofficesave: %s\n", err)
    }else {
      url, err := myjson.Get("url").String()
      if err == nil {
        os.MkdirAll(path.Join(cacheDir, filepath.Dir(fullpath)), 0755)
        abucket.Make(filepath.Dir(fullpath))
        localfile := path.Join(cacheDir,fullpath)
        errDownload := DownloadFile(localfile, url)
        if errDownload == nil {
          err := abucket.Upload(localfile, relativepath)
          if err != nil {
            fmt.Printf("----------- Onlyoffice error uploading [%s] [%s] 0000 [%s]\n", localfile, relativepath, err)
          }
        }
      }
    }
    c.JSON(map[string]int{"error": 0})
    c.Next()
  })

  app.Get("/new/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    target := c.URLParam("target")
    err:=abucket.Copy(relativepath, target, false)
    if err != nil {
      fmt.Printf("Error copying file %s -- %s \n", relativepath, target)
    }
    c.Redirect("/"+filepath.Dir(target))
    c.Next()
  })

  app.Get("/view/{path:path}", func (c iris.Context) {
    c.ViewData("BucketName", mybucket.Name())
    c.ViewData("Bucket", mybucket)
    c.ViewData("Path", c.Params().Get("path"))

    fmt.Printf("Getting path %s\n", c.Params().Get("path"))
    c.View("view.html")
    c.Next()
  })

  app.Get("/delete/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    username:= s.GetString("username")
    e2 := c.Values().Get("Enforcer").(*casbin.Enforcer)
    isAuthorized := false
    if e2 != nil {
      isAuthorized = e2.Enforce(username, relativepath, "delete")
    }
    if isAuthorized {
      fmt.Printf("Deleting path %s\n", relativepath)
      abucket.Delete(relativepath)
    }
    fmt.Printf("Redirecting to [%s]\n", filepath.Dir(fullpath))
    c.Redirect("/"+filepath.Dir(fullpath))
    c.Next()
  })

  app.Get("/jq/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)

    fullpath := c.Params().Get("path")

    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)

    expr := c.URLParam("expr")
    options := strings.Split(c.URLParam("options"), ",")
    targetFileset := c.URLParam("files")
    output := c.URLParam("output")
    translate := strings.Split(c.URLParam("translate"), ",")

    isDir := false
    tEls := []bucket.Element{}
    if targetFileset != "" {
      info, err := abucket.Stat(relativepath)
      if err == nil {
        if info != nil && info.IsDir() {
          isDir = true
        }
      }
      if isDir {
        els := abucket.List(relativepath)
        for _, el := range *els {
          m, _ := regexp.MatchString(targetFileset, el.Name())
          if m  {
            tEls = append(tEls, el)
          }
        }
        for _, l := range tEls {
          fmt.Printf("jq Selected elements %v\n", l.Name())
        }
      }
    }

    myfilepath := path.Join(cacheDir, filepath.Dir(fullpath))
    fmt.Printf("Creating directory %s\n", myfilepath)
    os.MkdirAll(myfilepath, 0755)
    command := ""
    allfiles := []string{}
    if isDir {
      for _, l := range tEls {
        fname := l.Name()
        err := abucket.Download(path.Join(relativepath,fname), path.Join(myfilepath, fname))
        if err != nil {
          fmt.Printf("jq download ---> %s\n", err)
        } else {
          allfiles = append(allfiles, path.Join(myfilepath, fname))
        }
      }
      command = fmt.Sprintf("./jq %s %s %s", options, expr, path.Join(myfilepath, targetFileset))
    } else {
      err := abucket.Download(relativepath, path.Join(cacheDir, fullpath))
      if err != nil {
        fmt.Printf("jq download ---> %s\n", err)
      }
      command = fmt.Sprintf("./jq %s %s %s", options, expr, path.Join(cacheDir, fullpath))
    }
    fmt.Printf("jq command ----> %s\n", command)
    var cmd *exec.Cmd
    if isDir {
      args := []string{}
      if len(options) == 0  {
        args = append(args, expr)
        args = append(args, allfiles...)
        cmd = exec.Command("jq", args...)
      } else {
        args = append(args, options...)
        args = append(args, expr)
        args = append(args, allfiles...)
        cmd = exec.Command("jq", args...)
      }
    } else {
      args := []string{}
      if len(options) == 0  {
        args = append(args, expr)
        args = append(args, path.Join(cacheDir, fullpath))
        cmd = exec.Command("jq", args...)
      } else {
        args = append(args, options...)
        args = append(args, expr)
        args = append(args, path.Join(cacheDir, fullpath))
        cmd = exec.Command("jq", args...)
      }
    }
    msg, err := cmd.CombinedOutput()
    if err != nil {
      fmt.Printf("jq exec ---> %s\n", err)
    }
    outputFilename := ""
    if output != "" {
      if isDir {
        outputFilename = path.Join(relativepath, output)
      } else {
        outputFilename = path.Join(filepath.Dir(relativepath), output)
      }
    }

    if len(translate)==2 {
      translated, err := gtranslate.TranslateWithParams(
        string(msg),
        gtranslate.TranslationParams{
          From: translate[0],
          To:   translate[1],
        },
      )
      if err != nil {
        fmt.Printf("Error in translation %s\n", err)
      }
      fmt.Printf("%s\n", string(translated))
      if output != "" {
        err = abucket.WriteStream(outputFilename, strings.NewReader(string(translated)), 0644)
        if err != nil {
          fmt.Printf("Error writing jq output %s\n", err)
        }
      }
      c.Writef(string(translated))
    } else {
      fmt.Printf("%s\n", string(msg))
      if output != "" {
        err = abucket.WriteStream(outputFilename, strings.NewReader(string(msg)), 0644)
        if err != nil {
          fmt.Printf("Error writing jq output %s\n", err)
        }
      }
      c.Writef(string(msg))
    }
  })

  defineGetActions(app, cacheDir)
  definePostActions(app, cacheDir)

  app.Get("/mkdir/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    fmt.Printf("Creating path %s\n", relativepath)

    err := abucket.Make(relativepath)
    if err != nil {
      fmt.Printf("Error creating path %s\n", err)
    }
    c.Next()
  })
  app.Get("/captcha/{captchaid:string}", func (c iris.Context) {
    captchaid:=c.Params().Get("captchaid")
    fmt.Printf("++++++++++++ %s\n", captchaid)
    captcha.WriteImage(c.ResponseWriter(), captchaid, 300, 200)
    c.Next()
  })
  app.Get("/login", func (c iris.Context) {
    s := mySessions.Start(c)
    c.ViewData("flasherror", s.GetFlashString("flasherror"))
    fullpath := c.URLParam("path")
    if fullpath != "" {
      c.ViewData("fullpath", fullpath)
    }
    c.ViewData(csrf.TemplateTag, csrf.TemplateField(c))
    c.ViewData("captcha", captcha.New())
    c.ViewData("fullpath", fullpath)
    c.View("login.html")
    c.Next()
  })
  app.Post("/login", func (c iris.Context) {
    email := c.FormValue("inputEmail")
    password := c.FormValue("inputPassword")
    fullpath := c.FormValue("fullpath")
    captchanum := c.FormValue("captcha")
    captchaid := c.FormValue("captchaid")
    fmt.Println(email, password, captchanum, captchaid)
    s := mySessions.Start(c)
    if !captcha.VerifyString(captchaid, captchanum) {
      s.SetFlash("flasherror", "Error in image number")
      if fullpath != "" {
        c.Redirect(fmt.Sprintf("/login?path=%s", fullpath))
      } else {
        c.Redirect("/login")
      }
      return
    }
    // see if path requested
    fmt.Printf("%%%%%%%%%%%%%%%%%%%%%%%%%%%% redirecting [%s]\n", fullpath)

    //set session values
    u := user.NewLocalUsers("users")
    defaultbucket := ""
    err :=u.CheckPassword(email, password)
    //if strings.HasSuffix(err.Error(), ": No such file or directory") {
    if  err != nil {
      // Not validated as a global user
      fmt.Printf("Validating user in %s\n", fullpath)
      defaultBucket := getDefaultBucket(s)
      abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
      fmt.Printf("Bucket using %+v\n", abucket)
      u := user.NewBucketUsers(abucket, relativepath)
      err :=u.CheckPassword(email, password)
      if err != nil {
        // See if we can validate as a bucket global user
        fmt.Printf("Validating user in bucket %s\n", abucket.Name())
        u := user.NewBucketUsers(abucket, "/")
        err :=u.CheckPassword(email, password)
        if err != nil {
          // login via ldap
          loggedIn := false
          email = strings.ToLower(email)
          if strings.HasSuffix(email, "imim.es") {
            ldapc, err := ldap.Dial("tcp", "172.20.4.10:389")
            if err != nil {
              fmt.Println(err)
            } else {
              defer ldapc.Close()
              fmt.Printf("U: %s, P: %s\n", email, password)
              err = ldapc.Bind(email, password)
              if err != nil {
                fmt.Println(err)
              } else {
                fmt.Printf("U: %s, P: %s -- Loged in\n", email, password)
                loggedIn = true
                s.Set("imimuser", "true")
              }
            }
          }

          if strings.HasSuffix(email, "parcdesalutmar.cat"){
            tmpUA := strings.Split(email, "@")
            targetUsername := tmpUA[0]
            if validatePSMAR(targetUsername, password) {
                fmt.Printf("U: %s, P: %s -- Loged in\n", email, password)
                loggedIn = true
                s.Set("imimuser", "true")
            }
          }
          if loggedIn == false {
            s.SetFlash("flasherror", "Error in email or password")
            if fullpath != "" {
              c.Redirect(fmt.Sprintf("/login?path=%s", fullpath))
            } else {
              c.Redirect("/login")
            }
            c.Next()
            return
          }
       }
      }
      defaultbucket = u.DefaultBucket(email)
      s.Set("RootPath", fullpath)
    } else {
      defaultbucket = u.DefaultBucket(email)
    }
    s.Set("name", "iris")
    s.Set("username", email)
    if defaultbucket  != "" {
      s.Set("default_bucket", defaultbucket)
    }
    s.Set("username", email)
    if fullpath != "" {
      c.Redirect(fullpath)
    } else  {
      c.Redirect("/")
    }
    c.Next()
  })

  app.Get("/logout", func (c iris.Context) {
    //set session values
    //s := mySessions.Start(c)
    //s.Delete("name")
    //s.Delete("username")
    mySessions.Destroy(c)
    c.Redirect("/login")
    c.Next()
  })

  app.Get("/data/{path:path}", func (c iris.Context) {
    fullpath := c.Params().Get("path")
    s := mySessions.Start(c)
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)

    reader, err := abucket.ReadStream(relativepath)
    if err != nil {
      fmt.Println(err)
      c.Writef(`{ "KO" }`)
    } else  {
      buf := new(bytes.Buffer)
      buf.ReadFrom(reader)
      s := buf.String() // Does a complete copy of the bytes in the buffer.
      c.Writef(s)
    }

    c.Next()
  })

  app.Get("/formfill/{path:path}", func (c iris.Context) {
    fmt.Printf("Starting formfill\n")
    s := mySessions.Start(c)
    c.ViewData(csrf.TemplateTag, csrf.TemplateField(c))
    token := csrf.Token(c)
    c.ViewData("CSRFToken", token)
    fullpath := c.Params().Get("path")
    decorate := c.URLParam("decorate")
    fmt.Printf("Decorating with %s\n", decorate)
    c.ViewData("Decorate", decorate)
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    extension := path.Ext(fullpath)
    c.ViewData("Path", fullpath)
    c.ViewData("DefaultBucket", defaultBucket)
    formpath := relativepath

    // if URL ends in .data
    if extension == ".data" {
      i := strings.Index(relativepath, ".formj")
      if i >= 0 {
        formpath=relativepath[:i]+".formj"
      } else {
        i := strings.Index(relativepath, ".form")
        if i >= 0 {
          formpath=relativepath[:i]+".form"
        }
      }
    }

    // form parameter needs to be fullpath
    if c.URLParam("form") != "" {
      formpath = c.URLParam("form")
    }

    rootTag :=  c.URLParam("root-tag")
    c.ViewData("rootTag", rootTag)

    // Remember formpath must be relative to bucket root
    fmt.Printf("Reading form: %s\n", formpath)
    reader, err := abucket.ReadStream(formpath)
    if err != nil {
      fmt.Printf("!!!! ERROR reading form: %s %s\n", formpath, err)
      c.ViewData("Components", "{ components: formioComponents } ")
    } else  {
      buf := new(bytes.Buffer)
      buf.ReadFrom(reader)
      s := buf.String() // Does a complete copy of the bytes in the buffer.
      re := regexp.MustCompile(`{{.CSRFToken}}`)
      s=re.ReplaceAllString(s, token)
      c.ViewData("Components", s)
      fmt.Printf("Form: %s read\n", formpath)
    }

    username:= s.GetString("username")
    c.ViewData("Values", fmt.Sprintf(`{ "username": "%s" }`, username))
    if extension != ".form" && extension != ".formj" {
      reader, err := abucket.ReadStream(relativepath)
      if err == nil {
        buf := new(bytes.Buffer)
        buf.ReadFrom(reader)
        ss := buf.String() // Does a complete copy of the bytes in the buffer.
        s1, _, _, errordata := jsonparser.Get([]byte(ss), "data")
        hasNoDataField := false
        if errordata != nil {
          if rootTag != "" {
            s1 = []byte(fmt.Sprintf(`{ "%s": %s }`, rootTag, ss))
          } else {
            s1 = []byte(ss)
          }
          hasNoDataField = true
        }
        myjson, err := simplejson.NewJson(s1)
        if err != nil {
          fmt.Printf("Simplejson error: %s\n", err)
        }else {
          myjson.Set("username", s.GetString("username"))
        }
        c.ViewData("Values", string(s1))
        s2 := ""
        if hasNoDataField {
          s2, _ = jsonparser.GetString([]byte(ss), "formname")
        }else {
          s2, _ = jsonparser.GetString([]byte(ss), "data", "formname")
        }
        if !strings.Contains(s2, "/") {
          s3 := s2
          s2 = path.Join(filepath.Dir(relativepath), s3)
          fmt.Printf("Switch to formname [%s]\n", s2)
        }
        fmt.Printf("formname: [%s] err [%s]\n", s2, err)
        if err == nil {
          reader, err := abucket.ReadStream(s2)
          if err == nil {
            fmt.Printf("Reading form %s\n", s2)
            buf := new(bytes.Buffer)
            buf.ReadFrom(reader)
            sdata := buf.String() // Does a complete copy of the bytes in the buffer.
            re := regexp.MustCompile(`{{.CSRFToken}}`)
            sdata=re.ReplaceAllString(sdata, token)
            c.ViewData("Components", sdata)
          } else {
            fmt.Printf("Unable to read form: [%s]\n", s2)
          }
        }
      }
    }

    fmt.Printf("Rendering %s\n", fullpath)
    if (filepath.Ext(fullpath) == ".formj") || (filepath.Ext(formpath) == ".formj") {
      fmt.Printf("Rendering formfillj\n")
      c.View("formfillj.html")
    } else if (filepath.Ext(fullpath) == ".formg") || (filepath.Ext(formpath) == ".formg") {
      c.View("grid.html")
    } else if (filepath.Ext(fullpath) == ".formcal") || (filepath.Ext(formpath) == ".formcal") {
      c.View("calendar.html")
    } else {
      c.View("formfill3.html")
    }
    c.Next()
  })

  app.Post("/formfill/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    rootTag :=  c.URLParam("root-tag")
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
    formdata, _, _, err := jsonparser.Get(rawData, "data")
    fmt.Printf("%+v\n", string(formdata))

    _, _, _, err = jsonparser.Get(rawData, "metadata")
    if err == nil {
      myjson, err := simplejson.NewJson(rawData)
      if err != nil {
        fmt.Printf("Simplejson error: %s\n", err)
      } else {
        myjson.Get("metadata").Set("username", s.GetString("username"))
        myjson.Get("metadata").Set("remoteip", c.Request().Header.Get("X-Forwarded-For"))
        myjson.Get("metadata").Set("realip", c.Request().Header.Get("X-Real-Ip"))
        fmt.Printf("Metadata %+v\n", myjson.Get("metadata"))
        var err2 error
        rawData, err2 = myjson.MarshalJSON()
        if err2 != nil {
          fmt.Printf("Simplejson marshall error: %s\n", err2)
        }
      }
    }

    filename, errfilename := jsonparser.GetString(rawData, "data", "filename")
    fmt.Printf("Filename: [%s] errfilename [%s]\n", filename, errfilename)
    extension := path.Ext(fullpath)

    formformat, _ := jsonparser.GetString(rawData, "data", "formformat")
    if formformat == "raw" {
      rawData = formdata
    }

    if rootTag != "" {
      formdata, _, _, err := jsonparser.Get(rawData, "data", rootTag)
      if err == nil {
        rawData = formdata
      }
    }

    myTargetFormFile := ""
    //if extension == ".data" {
    //no form and no filename, just save
    if extension != ".form" && extension != ".formj" && errfilename != nil {
      myTargetFormFile = relativepath
      err = abucket.WriteStream(relativepath, strings.NewReader(string(rawData)), 0644)
      if err != nil {
        fmt.Println(err)
      }
    } else {
      if errfilename != nil {
        u1 := uuid.NewV4()
        myTargetFormFile = relativepath+"."+u1.String()+".data"
        err = abucket.WriteStream(myTargetFormFile, strings.NewReader(string(rawData)), 0644)
        if err != nil {
          fmt.Println(err)
        }
      } else {
        mypath := filepath.Dir(relativepath)
        re, _ := regexp.Compile(`(\+counter:[^\+]+\+)`) 
        counter := re.FindString(filename)
        if counter != "" {
          i := ConsumeCounter(abucket, mypath, counter)
          filename = re.ReplaceAllString(filename, fmt.Sprintf("%d", i))
          myjson, err := simplejson.NewJson(rawData)
          if err != nil {
            fmt.Printf("Error parsing rawDat for filename %s\n", err)
          } else {
            myjson.Get("data").Set("filename", filename)
            var err2 error
            rawData, err2 = myjson.MarshalJSON()
            if err2 != nil {
              fmt.Printf("Simplejson marshall error: %s\n", err2)
            }
          }
        }
        myTargetFormFile=path.Join(mypath, filename)
        err = abucket.WriteStream(myTargetFormFile, strings.NewReader(string(rawData)), 0644)
        if err != nil {
          fmt.Println(err)
        }
      }
    }

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

    c.JSON(map[string]string{ "Result": "OK", "Redirect": redirect })

    c.Next()
  })


  app.Get("/uuid", func (c iris.Context) {
        u1 := uuid.NewV4()
        c.Writef(`{ "UUID": "%s" }`, u1.String())
  })

  app.Get("/form/{path:path}", func (c iris.Context) {
    s := mySessions.Start(c)
    c.ViewData("CSRFToken", csrf.Token(c))
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    c.ViewData("Path", fullpath)
    c.ViewData("DefaultBucket", defaultBucket)
    reader, err := abucket.ReadStream(relativepath)
    if err != nil {
      c.ViewData("Components", "{ components: formioComponents } ")
    }else {
      buf := new(bytes.Buffer)
      buf.ReadFrom(reader)
      s := buf.String() // Does a complete copy of the bytes in the buffer.
      c.ViewData("Components", s)
    }

    c.View("formbuilder3.html")
    c.Next()
  })

  app.Post("/form/{path:path}", iris.LimitRequestBodySize(10<<30), func (c iris.Context) {
    s := mySessions.Start(c)
    mypath := c.Params().Get("path")
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

    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(mypath, defaultBucket)
    err = abucket.WriteStream(relativepath, strings.NewReader(string(rawData)), 0644)
    if err != nil {
      fmt.Println(err)
    }
    c.Writef("OK")
    c.Next()
  })
  // Upload the file to the server
  // POST: http://localhost:8080/upload
  app.Post("/upload", iris.LimitRequestBodySize(10<<30), func(ctx iris.Context) {
    s := mySessions.Start(ctx)
    defaultBucket := getDefaultBucket(s)
    fpath := ctx.FormValue("fullPath")
    fullpath := ctx.FormValue("path")
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    // Get the file from the dropzone request

    file, info, err := ctx.FormFile("file")
    if err != nil {
      ctx.StatusCode(iris.StatusInternalServerError)
      ctx.Application().Logger().Warnf("Error while uploading: %v", err.Error())
      return
    }

    defer file.Close()
    //uploadsDir = ctx.FormValue("path")
    fmt.Printf("$$$$$$$ %s\n", uploadsDir)
    fname := info.Filename
    if fpath != "" {
      fname = fpath
      os.MkdirAll(path.Join(uploadsDir, fullpath, filepath.Dir(fpath)), 0755)
      abucket.Make(path.Join(relativepath, filepath.Dir(fpath)))

    }
    fmt.Printf(">>>>>>>>>> %+v\n", ctx.FormValues())

    // Create a file with the same name
    // assuming that you have a folder named 'uploads'
    /*
    out, err := os.OpenFile(path.Join(uploadsDir, fname),
    os.O_WRONLY|os.O_CREATE, 0666)

    if err != nil {
      ctx.StatusCode(iris.StatusInternalServerError)
      ctx.Application().Logger().Warnf("Error while preparing the new file: %v", err.Error())
      fmt.Println(err)
      return
    }
    defer out.Close()

    io.Copy(out, file)
    */
    fmt.Printf("------> Final file %s\n", path.Join(uploadsDir, fname))
    if fpath != "" {
      err = abucket.WriteStream(path.Join(relativepath, fpath), file, 0644)
    } else {
      err = abucket.WriteStream(path.Join(relativepath, fname), file, 0644)
    }
    if err != nil {
      fmt.Println(err)
    }
  })

  // Start SMTP
  go startSMTP()

  if konfig.String("certificate") != "" {
    certificate := konfig.String("certificate")
    key := konfig.String("privatekey")
    app.Run(iris.TLS(konfig.String("webserver"), certificate, key))
  } else {
    app.Run(iris.Addr(konfig.String("webserver")))
  }
}

/*
  app.Get("/set", func(ctx iris.Context) {

        //set session values
            s := mySessions.Start(ctx)
                s.Set("name", "iris")

                    //test if setted here
                        ctx.Writef("All ok session setted to: %s", s.GetString("name"))
                          })

                            app.Get("/get", func(ctx iris.Context) {
                                  // get a specific key, as string, if no found returns just an empty string
                                      s := mySessions.Start(ctx)
                                          name := s.GetString("name")

                                              ctx.Writef("The name on the /set was: %s", name)
                                                })

                                                  app.Get("/delete", func(ctx iris.Context) {
                                                        // delete a specific key
                                                            s := mySessions.Start(ctx)
                                                                s.Delete("name")
                                                                  })

                                                                    app.Get("/clear", func(ctx iris.Context) {
                                                                          // removes all entries
                                                                              mySessions.Start(ctx).Clear()
                                                                                })

                                                                                  app.Get("/update", func(ctx iris.Context) {
                                                                                        // updates expire date with a new date
                                                                                            mySessions.ShiftExpiration(ctx)
                                                                                              })

                                                                                                app.Get("/destroy", func(ctx iris.Context) {
                                                                                                      //destroy, removes the entire session data and cookie
                                                                                                          mySessions.Destroy(ctx)
                                                                                                            })
                                                                                                              // Note about destroy:
                                                                                                                //
                                                                                                                  // You can destroy a session outside of a handler too, using the:
                                                                                                                    // mySessions.DestroyByID
                                                                                                                    // mySessions.DestroyAll
                                                                                                                    */
