package main

import (
  "strconv"
  "github.com/casbin/casbin"
  "sort"
  "os"
  "fmt"
  "strings"
  "path"
  "path/filepath"
  "folduka/bucket"
  "encoding/json"
  "io/ioutil"
  "errors"
  "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
  "github.com/pdfcpu/pdfcpu/pkg/api"
  "os/exec"
  "github.com/kataras/iris/v12"
  "github.com/sethvargo/go-password/password"
)

func GetConfigFilename(abucket bucket.Bucket, targetPath, targetFilename string) string {
  policyFilename := ""

  curfilestat, err := abucket.Stat(targetPath)
  dir := "/"
  // if no err in filepath
  if err == nil {
    if curfilestat.IsDir() {
      dir = targetPath
    } else {
      dir = filepath.Dir(targetPath)
    }
  } else {
    // allow non existant targetPath
    //return policyFilename
  }

  elements := strings.Split(dir, "/")
  for i := len(elements) -1; i >=0 ; i-- {
    f, err := abucket.Stat(path.Join(path.Join(elements[0:i+1]...),"_config", targetFilename))
    if err == nil && f != nil {
      policyFilename =path.Join(path.Join(elements[0:i+1]...),"_config", targetFilename)
      break
    }
  }
  if policyFilename == "" {
    f, err := abucket.Stat(path.Join("_config", targetFilename))
    if err == nil && f != nil {
      policyFilename =path.Join("_config", targetFilename)
    }
  }
  return policyFilename
}

func GetPolicyFilename(abucket bucket.Bucket, dir string) string {
  policyFilename := ""
  elements := strings.Split(dir, "/")
  for i := len(elements) -1; i >=0 ; i-- {
    f, err := abucket.Stat(path.Join(path.Join(elements[0:i+1]...),"_config", "policy.json"))
    if err == nil && f != nil {
      policyFilename =path.Join(path.Join(elements[0:i+1]...),"_config", "policy.json")
      break
    }
  }
  if policyFilename == "" {
    f, err := abucket.Stat(path.Join("_config", "policy.json"))
    if err == nil && f != nil {
      policyFilename =path.Join("_config", "policy.json")
    }
  }
  return policyFilename
}

type Counters map[string]int

func ConsumeCounter(abucket bucket.Bucket, targetPath, name string) int {
  counterVal := 0

  curfilestat, err := abucket.Stat(targetPath)
  // if no err in filepath
  dir := "/"
  if err == nil {
    if curfilestat.IsDir() {
      dir = targetPath
    } else {
      dir = filepath.Dir(targetPath)
    }
  } else {
    return 0
  }

  counters := Counters{}

  cf := GetConfigFilename(abucket, dir, "counter.json")
  if cf == "" {
    cf = path.Join(dir, "_config", "counter.json")
  } else {
    b, _ := abucket.ReadStream(cf)
    bb, _ := ioutil.ReadAll(b)
    json.Unmarshal(bb, &counters)
  }

  fmt.Printf("Counter file :::::::::::::::::::::::::::::::: %s\n", cf)

  if val, ok := counters[name]; ok {
    counterVal = val
    counters[name]++
  } else {
    counters[name]=1
  }

  b2, _ := json.Marshal(counters)
  err = abucket.WriteStream(cf, strings.NewReader(string(b2)), 0644)
  if err != nil {
    fmt.Printf("Error writing counter: %s\n", err)
  }

  return counterVal
}

func OcrPdf(filename string) error {
  if filepath.Ext(filename) != ".pdf" {
    return errors.New("Not a pdf file")
  }
  tmpdir, _ := ioutil.TempDir(filepath.Dir(filename), "TMPDIR")
  listname, _ := SplitPdfInImages(filename, tmpdir)
  cmd := exec.Command("tesseract", listname, filename, "pdf")
  msg, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Printf("Error tessereact %s\n", msg)
    return err
  }
  return nil
}

func SplitPdfInImages(filename, dir string) (string, error) {
  //ºimagedir := path.Join(filepath.Dir(filename), dir)
  imagedir := dir
  os.MkdirAll(imagedir, 0755)
  config := pdfcpu.NewDefaultConfiguration()
  selectedPages := []string{"1-"}
  err := api.ExtractImagesFile(filename, imagedir, selectedPages, config)
  if err != nil {
    fmt.Printf("Error extracting images: %s\n", err)
  }

  files, _ := ioutil.ReadDir(dir)
  for _, f := range files {
    a:=strings.Split(f.Name(),"_")
    if len(a) == 3 {
      i, err := strconv.Atoi(a[1])
      if err == nil {
        newName:=fmt.Sprintf("%s_%03d_%s",a[0],i,a[2])
        os.Rename(path.Join(dir, f.Name()), path.Join(dir, newName))
      }
    }
  }

  files, _ = ioutil.ReadDir(dir)
  list := ""
  for _, f := range files {
    list = list + fmt.Sprintf("%s/%s\n", dir, f.Name())
  }

  tempfilelist, _ := ioutil.TempFile(filepath.Dir(filename), "LIST")
  ioutil.WriteFile(tempfilelist.Name(), []byte(list), 0644)

  return tempfilelist.Name(), nil
}

func defineUtilFunctions(app *iris.Application) {
  app.Get("/password", func (c iris.Context){
    // Generate a password that is 64 characters long with 10 digits, 10 symbols,
    // allowing upper and lower case letters, disallowing repeat characters.

    length := c.URLParamIntDefault("length", 8)
    digits := c.URLParamIntDefault("digits", 5)
    symbols := c.URLParamIntDefault("symbols", 0)
    bmixed, _ := c.URLParamBool("mixed")
    brepeat, _ := c.URLParamBool("repeat")
    res, err := password.Generate(length, digits, symbols, bmixed, brepeat)
    if err != nil {
      fmt.Printf("Error in password generation: %s\n", err)
    }
    c.Write([]byte(res))
  })

  app.Get("/dir/{path:path}", func(c iris.Context) {
    var els *[]bucket.Element
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    info, err := abucket.Stat(relativepath)
    if err != nil {
      fmt.Println(err)
    } else {
      if info != nil && info.IsDir() {
        // View folder
        c.ViewData("IsDir", true)
        actions, err := getActions(abucket, relativepath)
        if err != nil {
          fmt.Printf("Error in actions: %s\n", err)
        }
        c.ViewData("actions", actions)

        search := c.URLParam("search")
        if search != "" {
          fmt.Printf("Searching %s %s\n", search, relativepath)
          els = abucket.Search(search, relativepath)
        } else {
          els = abucket.List(relativepath)
        }
        sort.Sort(bucket.ElementByName(*els))
      }
    }
    type fi struct {
      Name string
      Size int64        // length in bytes for regular files; system-dependent for others
      IsDir bool
      Path string
    }
    myfi := []fi{}
    e := c.Values().Get("Enforcer").(*casbin.Enforcer)
    user := s.GetString("username")
    for _, a := range *els {
      resource := path.Join(a.Path(), a.Name())
      allow := e.Enforce(user, resource, "list")
      allow2 := e.Enforce(user, resource, "read")
      if allow || allow2 {
        myfi = append(myfi, fi{ Name: a.Name(), Size: a.Size(), IsDir: a.IsDir(), Path: resource })
      }
    }
    c.JSON(myfi)
    c.Next()
  })
}
