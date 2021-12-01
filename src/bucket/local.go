package bucket

import (
  "strings"
  "path/filepath"
  "os"
  "fmt"
  "io"
  "path"
  "io/ioutil"
  "bytes"
  "time"
  pt "github.com/monochromegane/the_platinum_searcher"
)


type LocalElement struct {
  file os.FileInfo
  path string
  prefix string
  bucket Bucket
}

func (e LocalElement) Size() int64 {
  return e.file.Size()
}
func (e LocalElement) Prefix() string {
  return e.prefix
}

func (e LocalElement) Bucket() Bucket {
  return e.bucket
}

func (e LocalElement) Path() string {
  return e.path
}

func (e LocalElement) Name() string {
  return e.file.Name()
}
func (e LocalElement) ModTime() time.Time {
  return e.file.ModTime()
}

func (e LocalElement) IsDir() bool {
  return e.file.IsDir()
}

var (
//  folder_icon_tags = `<span class="oi oi-icon-folder" title="folder" aria-hidden="true"></span>`
//  file_icon_tags = `<span class="oi oi-icon-file" title="file" aria-hidden="true"></span>`
  folder_icon_tags = `<img src="/svg/folder.svg" width="16" alt="folder">`
  file_icon_tags = `<img src="/svg/file.svg" width="16" alt="file">`
)

func (e LocalElement) Display() string {
  if e.file.IsDir() {
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
  /*
    //fmt.Printf("########### %s ----- %s ---- %s\n", e.path, e.file.Name(), path.Join("/", e.path, e.file.Name()))
    if e.prefix != "" {
      return fmt.Sprintf("<tr><td>%s&nbsp;<a href=\""+path.Join("/",e.prefix, e.path, e.file.Name())+"\">" + e.file.Name() + "</a></td></tr>", folder_icon_tags)
    } else {
      return fmt.Sprintf("<tr><td>%s&nbsp;<a href=\""+path.Join("/",e.path, e.file.Name())+"\">" + e.file.Name() + "</a></td></tr>", folder_icon_tags)
    }
  } else {
    if e.prefix != "" {
      return fmt.Sprintf("<tr><td>%s&nbsp;<a href=\""+path.Join("/download", "/",e.prefix, e.path, e.file.Name())+"\">" + e.file.Name() + "</a></td></tr>", file_icon_tags)
    } else {
      return fmt.Sprintf("<tr><td>%s&nbsp;<a href=\""+path.Join("/download", "/",e.path, e.file.Name())+"\">" + e.file.Name() + "</a></td></tr>", file_icon_tags)
    }
  }
    */
}

type LocalBucket struct {
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

func NewLocalBucket(name, root, host, user, password string) *LocalBucket {
  os.MkdirAll(root, 0755)
  return &LocalBucket{ name: name, root: root, host: host, user: user, password: password }
}

func (b LocalBucket)Name() string {
  return b.name
}

func (b LocalBucket)Prefix() string {
  return b.prefix
}

func (b *LocalBucket)SetPrefix(prefix string) {
  b.prefix = prefix
}
func (b LocalBucket)GetRoot() string {
  return b.root
}

func (b LocalBucket)DisplayHeader() string {
  return `<table class="table table-striped">`
}

func (b LocalBucket)DisplayFooter() string {
  return "</table>"
}

func (b *LocalBucket)ListRoot() *[]Element {
  if b.Elements != nil {
    return b.Elements
  }
  els := []Element{}

  files, err := ioutil.ReadDir(b.root)
  if err != nil {
    fmt.Println(b.root)
    return &els
  }
  for _, f := range files {
    els = append(els, LocalElement{ file: f, path: "", prefix: b.prefix, bucket: b })
  }
  b.Elements = &els
  return &els
}


func (b *LocalBucket)List(mypath string) *[]Element {
  els := []Element{}

  files, err := ioutil.ReadDir(path.Join(b.root, mypath))
  if err != nil {
    fmt.Println(b.root)
    return &els
  }
  for _, f := range files {
    els = append(els, LocalElement{ file: f, path: mypath, prefix: b.prefix, bucket: b })
  }
  b.Elements = &els
  return &els
}

func (b *LocalBucket)Delete(mypath string) error {
  err := os.Remove(path.Join(b.root, mypath))
  return err
}

func (b *LocalBucket)Make(mypath string) error {
  err := os.MkdirAll(path.Join(b.root, mypath), 0755)
  return err
}

func (b *LocalBucket)Stat(mypath string) (os.FileInfo, error) {
  if mypath == "" {
    mypath = "./"
  }
  return os.Stat(path.Join(b.root, mypath))
}

func (b *LocalBucket)WriteStream(mypath string, stream io.Reader, fm os.FileMode) error {
  f, err := os.Create(path.Join(b.root, mypath))
  if err != nil {
    return err
  }
  defer f.Close()
  _, err = io.Copy(f, stream)
  return err
}

func (b *LocalBucket)ReadStream(mypath string) (io.Reader, error) {
  dat, err := ioutil.ReadFile(path.Join(b.root, mypath))
  if err != nil {
    return nil, err
  }
  return bytes.NewReader(dat), nil
}

func (b *LocalBucket)Download(mypath, local string) error {
  _, err := copy(path.Join(b.root, mypath), local)
  return err
}

func (b *LocalBucket)Upload(local, mypath string) error {
  _, err := copy(local, path.Join(b.root, mypath))
  return err
}

func copy(src, dst string) (int64, error) {
  sourceFileStat, err := os.Stat(src)
  if err != nil {
    return 0, err
  }

  if !sourceFileStat.Mode().IsRegular() {
    return 0, fmt.Errorf("%s is not a regular file", src)
  }

  source, err := os.Open(src)
  if err != nil {
    return 0, err
  }
  defer source.Close()

  destination, err := os.Create(dst)
  if err != nil {
    return 0, err
  }
  defer destination.Close()
  nBytes, err := io.Copy(destination, source)
  return nBytes, err
}

func (b *LocalBucket) AddAction(key, value string) {
  b.Actions[key]=value
}

func (b *LocalBucket) DelAction(key string){
  delete(b.Actions, key)
}

func (b *LocalBucket) AddElementAction(key, value string){
  b.ElementActions[key]=value
}

func (b *LocalBucket) DelElementAction(key string){
  delete(b.ElementActions, key)
}

func (b *LocalBucket)  Rename(source, destination string, isoverwrite bool) error {
  return os.Rename(path.Join(b.root, source), path.Join(b.root, destination))
}

func (b *LocalBucket)  Copy(source, destination string, isoverwrite bool) error {
  _, err := mycopy(path.Join(b.root, source), path.Join(b.root, destination))
  return err
}

func mycopy(src, dst string) (int64, error) {
  sourceFileStat, err := os.Stat(src)
  if err != nil {
    return 0, err
  }

  if !sourceFileStat.Mode().IsRegular() {
    return 0, fmt.Errorf("%s is not a regular file", src)
  }

  source, err := os.Open(src)
  if err != nil {
    return 0, err
  }
  defer source.Close()

  destination, err := os.Create(dst)
  if err != nil {
    return 0, err
  }
  defer destination.Close()
  nBytes, err := io.Copy(destination, source)
  return nBytes, err
}

func (b *LocalBucket) Search(searchPattern, mypath string) *[]Element {
  var tpl bytes.Buffer
  var tpl2 bytes.Buffer
  pt := pt.PlatinumSearcher{Out: &tpl, Err: &tpl2}
  mmmpath := path.Join(b.root, mypath)
  fmt.Printf("Searching in path %s\n", mmmpath)
  exitCode := pt.Run([]string{ "-l", "--nocolor", "-i", searchPattern, mmmpath })
  fmt.Println(exitCode)
  fmt.Println(tpl.String())
  fmt.Println(tpl2.String())
  files := strings.Split(tpl.String(), "\n")
  els := []Element{}
  for _, fn := range files {
    fmt.Println(fn)
    f, _ := os.Stat(fn)
    if f != nil {
      els = append(els, LocalElement{ file: f, path: filepath.Dir(fn), prefix: b.Prefix(), bucket: b })
    }

  }
  return &els
}
