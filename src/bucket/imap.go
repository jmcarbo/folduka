package bucket

import (
  "os"
  "fmt"
  "io"
  "path"
  "io/ioutil"
  "bytes"
  "time"
  //"strings"
//  "path/filepath"
  "github.com/emersion/go-message/charset"
  "golang.org/x/text/encoding/charmap"
  "github.com/emersion/go-imap"
  "github.com/emersion/go-imap/client"
  "github.com/emersion/go-message/mail"
  "crypto/sha256"
)

func init() {
  charset.RegisterEncoding("iso-8859-7", charmap.ISO8859_7)
  charset.RegisterEncoding("iso-8859-8", charmap.ISO8859_8)
  charset.RegisterEncoding("windows-1253", charmap.Windows1253)
  charset.RegisterEncoding("windows-1254", charmap.Windows1254)
  charset.RegisterEncoding("windows-1257", charmap.Windows1257)
  charset.RegisterEncoding("windows-1258", charmap.Windows1258)
  charset.RegisterEncoding("windows-1255", charmap.Windows1255)
  charset.RegisterEncoding("windows-1256", charmap.Windows1256)
}

type ImapElement struct {
  file os.FileInfo
  path string
  prefix string
  bucket Bucket
}

func (e ImapElement) Size() int64 {
  return e.file.Size()
}
func (e ImapElement) Prefix() string {
  return e.prefix
}

func (e ImapElement) Bucket() Bucket {
  return e.bucket
}

func (e ImapElement) Path() string {
  return e.path
}

func (e ImapElement) Name() string {
  return e.file.Name()
}
func (e ImapElement) ModTime() time.Time {
  return e.file.ModTime()
}

func (e ImapElement) IsDir() bool {
  return e.file.IsDir()
}

var (
//  folder_icon_tags = `<span class="oi oi-icon-folder" title="folder" aria-hidden="true"></span>`
//  file_icon_tags = `<span class="oi oi-icon-file" title="file" aria-hidden="true"></span>`
)

func (e ImapElement) Display() string {
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

type ImapBucket struct {
  host string
  user string
  password string
  name string
  root string
  prefix string

  Elements *[]Element

  Actions map[string]string
  ElementActions map[string]string

  cl *client.Client
}

func NewImapBucket(name, root, host, user, password string) *ImapBucket {
  os.MkdirAll(root, 0755)
  return &ImapBucket{ name: name, root: root, host: host, user: user, password: password }
}

func (b *ImapBucket) LoginImap() *client.Client {
  var err2 error
  b.cl, err2 = client.DialTLS(b.host, nil)
  if err2 != nil {
    fmt.Printf("IMAP LOGIN %s\n", err2)
    return nil
  }
  //defer c2.Logout()
  fmt.Println("Connected")
  // Login
  if err := b.cl.Login(b.user, b.password); err != nil {
    fmt.Println(err)
    return nil
  }
  fmt.Println("Logged in")
  return b.cl
}

func (b *ImapBucket)LogoutImap() {
  b.cl.Logout()
}

func (b ImapBucket)Name() string {
  return b.name
}

func (b ImapBucket)Prefix() string {
  return b.prefix
}

func (b *ImapBucket)SetPrefix(prefix string) {
  b.prefix = prefix
}
func (b ImapBucket)GetRoot() string {
  return b.root
}

func (b ImapBucket)DisplayHeader() string {
  return `<table class="table table-striped">`
}

func (b ImapBucket)DisplayFooter() string {
  return "</table>"
}

func (b *ImapBucket)ListRoot() *[]Element {
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


func (b *ImapBucket)SaveMessage(folder string, msg *imap.Message) error {
  //log.Printf("* %+v", msg)
  //log.Printf("* " + msg.Envelope.Subject)
  //log.Printf("* %+v", msg.Body)
  //log.Printf("* %+v", msg.Envelope)
  //log.Printf("* %+v", msg.Body)
  //messageId := base64.URLEncoding.EncodeToString([]byte(msg.Envelope.MessageId))
  targetfolder := b.root
  if folder != "" {
    targetfolder = path.Join(targetfolder, folder)
  }
  /*
  fmt.Printf("Saving message %s\n", string(msg.Envelope.MessageId))
  if _, err := os.Stat(path.Join(targetfolder, messageId)); err == nil {
    log.Printf("Skipping messsage %s\n", messageId)
    return nil
  }
  */

  var mmm string
  for _,v := range msg.Body {
    //log.Printf("---> %+v\n", v)
    mmm = fmt.Sprintf("%s", v)
    break
  }

  sum := sha256.Sum256([]byte(mmm))
  messageId := fmt.Sprintf("%x", sum)
  err := ioutil.WriteFile(path.Join(targetfolder, messageId + ".email"), []byte(mmm), 0644)
  if err != nil {
    fmt.Println(err)
    return err
  }
  return nil
}

func (b *ImapBucket)List(mypath string) *[]Element {
  els := []Element{}

  if mypath == "" {
    c:=b.LoginImap()
    defer b.LogoutImap()
    // Select INBOX
    mbox, err := c.Select("INBOX", false)
    if err != nil {
      fmt.Println(err)
      return nil
    }
    fmt.Println("Flags for INBOX:", mbox.Flags)

    // Get the last 4 messages
    from := uint32(1)
    to := mbox.Messages
    /*
    if mbox.Messages > 3 {
      // We're using unsigned integers here, only substract if the result is > 0
      from = mbox.Messages - 3
    }
    */
    seqset := new(imap.SeqSet)
    seqset.AddRange(from, to)

    files, err := ioutil.ReadDir(path.Join(b.root, mypath))
    if err != nil {
      fmt.Println(b.root)
    } else {
      for _, f := range files {
        err2 := os.Remove(path.Join(b.root, mypath, f.Name()))
        if err2 != nil {
          fmt.Println(err)
        }
      }
    }

    messages := make(chan *imap.Message, 10)
    done := make(chan error, 1)
    go func() {
      section := &imap.BodySectionName{}
      items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope}

      done <- c.Fetch(seqset, items, messages)
    }()

    fmt.Println("Last 4 messages:")
    for msg := range messages {
      fmt.Println("* " + msg.Envelope.Subject)
      err := b.SaveMessage(mypath, msg)
      if err != nil {
        fmt.Printf("Erro saving imap message %s\n", err)
      }
    }

    if err := <-done; err != nil {
      fmt.Println(err)
      return nil
    }

    fmt.Println("Done!")
  }

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

func (b *ImapBucket)Delete(mypath string) error {
  err := os.Remove(path.Join(b.root, mypath))
  return err
}

func (b *ImapBucket)Make(mypath string) error {
  err := os.MkdirAll(path.Join(b.root, mypath), 0755)
  return err
}

func (b *ImapBucket)Stat(mypath string) (os.FileInfo, error) {
  if mypath == "" {
    mypath = "./"
  }
  return os.Stat(path.Join(b.root, mypath))
}

func (b *ImapBucket)WriteStream(mypath string, stream io.Reader, fm os.FileMode) error {
  f, err := os.Create(path.Join(b.root, mypath))
  if err != nil {
    return err
  }
  defer f.Close()
  _, err = io.Copy(f, stream)
  return err
}

func (b *ImapBucket)ReadStream(mypath string) (io.Reader, error) {
  dat, err := ioutil.ReadFile(path.Join(b.root, mypath))
  if err != nil {
    return nil, err
  }
  return bytes.NewReader(dat), nil
}

func (b *ImapBucket)Download(mypath, local string) error {
  _, err := copy(path.Join(b.root, mypath), local)
  return err
}

func (b *ImapBucket)Upload(local, mypath string) error {
  _, err := copy(local, path.Join(b.root, mypath))
  return err
}

func (b *ImapBucket) AddAction(key, value string) {
  b.Actions[key]=value
}

func (b *ImapBucket) DelAction(key string){
  delete(b.Actions, key)
}

func (b *ImapBucket) AddElementAction(key, value string){
  b.ElementActions[key]=value
}

func (b *ImapBucket) DelElementAction(key string){
  delete(b.ElementActions, key)
}

func (b *ImapBucket)  Rename(source, destination string, isoverwrite bool) error {
  return os.Rename(path.Join(b.root, source), path.Join(b.root, destination))
}

func (b *ImapBucket)  Copy(source, destination string, isoverwrite bool) error {
  _, err := mycopy(path.Join(b.root, source), path.Join(b.root, destination))
  return err
}

func ParseMessageDate(r io.Reader) time.Time {
 m, err := mail.CreateReader(r)
 if err != nil {
     fmt.Println("Reading Message Date")
     return time.Time{}
 }

  dateString, _ := m.Header.Text("Date")
  fmt.Printf("%s\n", dateString)
  if len(dateString) == 0 {
    return time.Now() 
  }
  //dateStringArray := strings.Split(m.Header.Header["Date"][0], " boundary")
  //dateString := dateStringArray[0]
  const longForm = "Jan 2, 2006 at 3:04pm (MST)"
// "15 Nov 2016 13:12:06 +0000"
  const longForm2 = "02 Jan 2006 15:04:05 -0700"
  const longForm3 = "Mon, 2 Jan 2006 15:04:05 -0700 (MST)"
  const longForm4 = "Mon, 2 Jan 2006 15:04:05 -0700"
  const longForm5 = "Mon, 02 Jan 2006 15:04:05 MST"
  const longForm7 = "Mon, 02 Jan 2006 15:04 MST"
  const longForm6 = "Mon, 02 Jan 2006 15:04:05 UT"
  const longForm8 = "2 Jan 2006 15:04:05 -0700"
  // "Sat, 14 Apr 2018 16:31:20 +0200"
  //t, err := time.Parse(time.RFC3339, m.Header.Header["Date"][0])
  var t time.Time
  t, err = time.Parse(time.RFC1123Z, dateString)
   if err != nil {
     t, err = time.Parse(time.RFC822Z, dateString)
     if err != nil {
       t, err = time.Parse(longForm2, dateString)
       if err != nil {
         t, err = time.Parse(longForm3, dateString)
         if err != nil {
           t, err = time.Parse(longForm4, dateString)
           if err != nil {
             t, err = time.Parse(longForm5, dateString)
             if err != nil {
               t, err = time.Parse(longForm6, dateString)
               if err != nil {
                 t, err = time.Parse(longForm7, dateString)
                 if err != nil {
                   t, err = time.Parse(longForm8, dateString)
                   if err != nil {
                     fmt.Println(err)
		     t = time.Now()
                   }
                 }
               }
             }
           }
         }
       }
     }
   }
   fmt.Println(t)
/* 
     for {
         p, err := m.NextPart()
         if err == io.EOF {
             break
         } else if err != nil {
             log.Fatal(err)
         }
 
         //t, _, _ := p.Header.ContentType()
         //log.Println("A part with type", t)
         log.Printf("%+v\n", p.Header)
     }
*/
  return t
}

func (b *ImapBucket) Search(searchPattern, mypath string) *[]Element {
  return nil
}

