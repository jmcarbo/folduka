package bucket

import (
  "github.com/studio-b12/gowebdav"
  "os"
  "fmt"
  "io"
  "net/http"
  "crypto/tls"
  "path"
  "io/ioutil"
  "strings"
  "time"
)


type DAVElement struct {
  file os.FileInfo
  path string
  prefix string
  bucket Bucket
}

func (e DAVElement) Size() int64 {
  return e.file.Size()
}
func (e DAVElement) Name() string {
  return e.file.Name()
}

func (e DAVElement) Prefix() string {
  return e.prefix
}

func (e DAVElement) Bucket() Bucket {
  return e.bucket
}

func (e DAVElement) ModTime() time.Time {
  return e.file.ModTime()
}
func (e DAVElement) IsDir() bool {
  return e.file.IsDir()
}

func (e DAVElement) Path() string {
  return e.path
}

func (e DAVElement) Display() string {
  /*
  if e.file.IsDir() {
    //fmt.Printf("########### %s ----- %s ---- %s\n", e.path, e.file.Name(), path.Join("/", e.path, e.file.Name()))
    return "<tr><td><a href=\""+path.Join("/",e.path, e.file.Name())+"\">" + e.file.Name() + "</a></td></tr>"
  } else {
    return "<tr><td><a href=\""+path.Join("/download", "/",e.path, e.file.Name())+"\">" + e.file.Name() + "</a></td></tr>"
  }
  */
  if e.file.IsDir() {
    //fmt.Printf("########### %s ----- %s ---- %s\n", e.path, e.file.Name(), path.Join("/", e.path, e.file.Name()))
    if e.prefix != "" {
      return fmt.Sprintf("<tr><td>%s&nbsp;<a href=\""+path.Join("/",e.prefix, e.path, e.file.Name())+"\">" + e.file.Name() + "</a></td><td></td><td>%s</td></tr>", folder_icon_tags, e.file.ModTime().Format("2006-01-02 15:04:05"))
    } else {
      return fmt.Sprintf("<tr><td>%s&nbsp;<a href=\""+path.Join("/",e.path, e.file.Name())+"\">" + e.file.Name() + "</a></td><td></td><td>%s</td></tr>", folder_icon_tags, e.file.ModTime().Format("2006-01-02 15:04:05"))
    }
  } else {
    if e.prefix != "" {
      return fmt.Sprintf("<tr><td>%s&nbsp;<a href=\""+path.Join("/download", "/",e.prefix, e.path, e.file.Name())+"\">" + e.file.Name() + "</a></td><td>%-5d KB</td><td>%s</td></tr>", file_icon_tags, e.file.Size()/1000, e.file.ModTime().Format("2006-01-02 15:04:05"))
    } else {
      return fmt.Sprintf("<tr><td>%s&nbsp;<a href=\""+path.Join("/download", "/",e.path, e.file.Name())+"\">" + e.file.Name() + "</a></td><td>%-5d KB</td><td>%s</td></tr>", file_icon_tags, e.file.Size()/1000, e.file.ModTime().Format("2006-01-02 15:04:05"))
    }
  }
}

type DAVBucket struct {
  host string
  user string
  password string
  name string
  root string
  prefix string

  Elements *[]Element

  Actions map[string]string
  ElementActions map[string]string
}

func (b *DAVBucket)SetPrefix(prefix string) {
  b.prefix = prefix
}

func NewDAVBucket(name, root, host, user, password string) *DAVBucket {
  http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
  return &DAVBucket{ name: name, root: root, host: host, user: user, password: password }
}

func (b DAVBucket)Name() string {
  return b.name
}

func (b DAVBucket)Prefix() string {
  return b.prefix
}

func (b DAVBucket)GetRoot() string {
  return "/"
}

func (b DAVBucket)DisplayHeader() string {
  return `<table class="table table-striped">`
}

func (b DAVBucket)DisplayFooter() string {
  return "</table>"
}

func (b *DAVBucket)ListRoot() *[]Element {
  if b.Elements != nil {
    return b.Elements
  }
  els := []Element{}

  c := gowebdav.NewClient(b.root, b.user, b.password)
  err := c.Connect()
  if err != nil {
    fmt.Println(err)
    return &els
  }
  fmt.Printf("Connection %+v\n", c)
  files, err := c.ReadDir("/")
  if err != nil {
    fmt.Println(b.root)
    return &els
  }
  for _, f := range files {
    if !strings.HasPrefix(path.Base(f.Name()), ".") && !strings.HasPrefix(path.Base(f.Name()), "~") {
      els = append(els, DAVElement{ file: f, path: "/", prefix: b.prefix, bucket: b })
    }
  }
  b.Elements = &els
  return &els
}


func (b *DAVBucket)List(mypath string) *[]Element {
  /*
  if b.Elements != nil {
    return b.Elements
  }
  */
  els := []Element{}

  c := gowebdav.NewClient(b.root, b.user, b.password)
  err := c.Connect()
  if err != nil {
    fmt.Println(err)
    return &els
  }
  fmt.Printf("Connection %+v\n", c)
  if mypath == "" {
    mypath = "/"
  }
  files, err := c.ReadDir(mypath)
  if err != nil {
    fmt.Println(b.root)
    return &els
  }
  for _, f := range files {
    if !strings.HasPrefix(path.Base(f.Name()), ".") && !strings.HasPrefix(path.Base(f.Name()), "~") {
      els = append(els, DAVElement{ file: f, path: mypath, prefix: b.prefix, bucket: b })
    }
  }
  b.Elements = &els
  return &els
}

func (b *DAVBucket)Delete(path string) error {
  c := gowebdav.NewClient(b.root, b.user, b.password)
  err := c.Connect()
  if err != nil {
    fmt.Println(err)
    return err
  }
  err = c.Remove(path)
  return err
}
func (b *DAVBucket)Make(path string) error {
  c := gowebdav.NewClient(b.root, b.user, b.password)
  err := c.Connect()
  if err != nil {
    fmt.Println(err)
    return err
  }
  err = c.MkdirAll(path, 0755)
  return err
}

func (b *DAVBucket)Stat(path string) (os.FileInfo, error) {
  c := gowebdav.NewClient(b.root, b.user, b.password)
  err := c.Connect()
  if err != nil {
    fmt.Println(err)
    return gowebdav.File{}, err
  }
  return c.Stat(path)
}

func (b *DAVBucket)WriteStream(path string, stream io.Reader, fm os.FileMode) error {
  c := gowebdav.NewClient(b.root, b.user, b.password)
  err := c.Connect()
  if err != nil {
    fmt.Println(err)
    return err
  }
  return c.WriteStream(path, stream, fm)
}

func (b *DAVBucket)ReadStream(path string) (io.Reader, error) {
  c := gowebdav.NewClient(b.root, b.user, b.password)
  err := c.Connect()
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  return c.ReadStream(path)
}

func (b *DAVBucket)Download(path, local string) error {
  c := gowebdav.NewClient(b.root, b.user, b.password)
  err := c.Connect()
  if err != nil {
    fmt.Println(err)
    return err
  }

  reader, err := c.ReadStream(path)
  if err != nil {
    return err
  }

  file, err := os.Create(local)
  if err != nil {
    return err
  }
  defer file.Close()

  _, err = io.Copy(file, reader)
  return err
}

func (b *DAVBucket)Upload(local, path string) error {
  c := gowebdav.NewClient(b.root, b.user, b.password)
  err := c.Connect()
  if err != nil {
    fmt.Println(err)
    return err
  }

  bytes, _ := ioutil.ReadFile(local)
  err = c.Write(path, bytes, 0644)
  return err
}

func (b *DAVBucket) AddAction(key, value string) {
  b.Actions[key]=value
}

func (b *DAVBucket) DelAction(key string){
  delete(b.Actions, key)
}

func (b *DAVBucket) AddElementAction(key, value string){
  b.ElementActions[key]=value
}

func (b *DAVBucket) DelElementAction(key string){
  delete(b.ElementActions, key)
}

func (b *DAVBucket)  Rename(source, destination string, isoverwrite bool) error {
  c := gowebdav.NewClient(b.root, b.user, b.password)
  err := c.Connect()
  if err != nil {
    fmt.Println(err)
    return err
  }
  return c.Rename(source, destination, isoverwrite)
}

func (b *DAVBucket)  Copy(source, destination string, isoverwrite bool) error {
  c := gowebdav.NewClient(b.root, b.user, b.password)
  err := c.Connect()
  if err != nil {
    fmt.Println(err)
    return err
  }
  return c.Copy(source, destination, isoverwrite)
}

func (b *DAVBucket) Search(searchPattern, mypath string) *[]Element {
  return nil
}

