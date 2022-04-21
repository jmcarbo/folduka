package main

import (
	"log"
	"time"
  "os"
  "errors"
  "net/mail"
  "io"
  "io/ioutil"
  "crypto/tls"
  "fmt"
  "github.com/emersion/go-smtp"
  "github.com/satori/go.uuid"
  "strings"
  "github.com/lalamove/konfig"
  "github.com/jhillyerd/enmime"
  "path"
  "path/filepath"
  "encoding/json"
  "github.com/kataras/iris"
)

type Backend struct{}

func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
  fmt.Printf("Loging in with %s -- %s\n", username, password)
  if username != "6imim" || password != "j(6}@zCh)MA]utYLnJ1Rb" {
    fmt.Printf("Loging FAILURE in with %s -- %s\n", username, password)
    return nil, errors.New("Invalid username or password")
  } else {
    fmt.Printf("Loging SUCCESS in with %s -- %s\n", username, password)
  }
  return &Session{ IsAnonymous: false, Username: username, Password: password }, nil
}

// Require clients to authenticate using SMTP AUTH before sending emails
func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	//return nil, smtp.ErrAuthRequired
        return &Session{ IsAnonymous: true }, nil
}


type smtpFolderConfig struct {
  ExplodeAttachments bool
  OCRPdf bool
}

// A Session is returned after successful login.
type Session struct{
  Username string
  Password string
  IsAnonymous bool
  From string
  To []string
}

func (s *Session) Data(r io.Reader) error {
	from := s.From
	to := s.To
  cacheDir := "./cache"

	fmt.Println("Receiving message:", from, to)
/*
        if !u.IsAnonymous {
          forwardMail(to, r)
        } else {
          fmt.Printf("To ---> %v\n", to)
          forwardMail([]string{"jmcarbo@imim.es"}, r)
        }
        */

  if b, err := ioutil.ReadAll(r); err != nil {
    return err
  } else {
    go func() {
      fmt.Println("Data:", string(b))
      defaultBucket := konfig.String("default_bucket")
      domain := konfig.String("domain")
      fmt.Printf("Domain <<<<<<<<<<<< %s\n", domain)
      for _, ito := range to {
        tto := strings.TrimSuffix(ito, "@"+domain)
        tto = strings.Replace(tto, "+", "/", -1)
        abucket, relativepath := LoadBucketFromPath(tto, defaultBucket)
        cfn := GetConfigFilename(abucket, relativepath, "smtp.json")
        if cfn != "" {
          smtpconfig := smtpFolderConfig{}
          bb, err :=abucket.ReadStream(cfn)
          if err != nil {
            fmt.Println(err)
          } else {
            smtpbytes, err := ioutil.ReadAll(bb)
            if err != nil {
              fmt.Println(err)
            } else {
              if err := json.Unmarshal([]byte(smtpbytes), &smtpconfig); err != nil {
                fmt.Printf("display parsing smtp.json error %s\n", err)
              }
            }
          }
          // Parse message body with enmime.
          env, err := enmime.ReadEnvelope(strings.NewReader(string(b)))
          if err != nil {
                fmt.Print(err)
                return
          }
          // enmime can decode quoted-printable headers.
          fmt.Printf("Subject: %v\n", env.GetHeader("Subject"))
          subject :=  env.GetHeader("Subject")
          myuuid := uuid.NewV4()
          mypath := relativepath
          myfilename := path.Join(relativepath, myuuid.String()+".mail")
          if !strings.HasPrefix(subject, "Message") {
            abucket.Make(path.Join(relativepath, subject))
            mypath = path.Join(relativepath, subject)
            myfilename = path.Join(relativepath, subject, myuuid.String()+".mail")
          }
          err = abucket.WriteStream(myfilename, strings.NewReader(string(b)), 0644)
          if err != nil {
            fmt.Printf("Error writing smtp message %s\n", err)
          }
          fmt.Printf("SMTPCONFIG %+v\n", smtpconfig)
          if smtpconfig.ExplodeAttachments {
            for _, a := range env.Attachments {
              extension := filepath.Ext(a.FileName)
              afilename := path.Join(mypath, strings.TrimSuffix(a.FileName, extension) + myuuid.String() + extension)
              if extension == ".pdf" {
                localfilename :=path.Join(cacheDir, afilename)
                os.MkdirAll(filepath.Dir(localfilename), 0755)
                ioutil.WriteFile(localfilename, a.Content, 0644)
                if smtpconfig.OCRPdf {
                  PutQueue(localfilename+";"+afilename+";"+abucket.Prefix())
                } else {
                  abucket.Upload(localfilename, afilename)
                }
              } else {
                err = abucket.WriteStream(afilename, strings.NewReader(string(a.Content)), 0644)
              }
              if err != nil {
                fmt.Printf("Error writing smtp attachment %s\n", err)
              }
            }
            for _, a := range env.Inlines {
              extension := filepath.Ext(a.FileName)
              afilename := path.Join(mypath, strings.TrimSuffix(a.FileName, extension) + myuuid.String() + extension)
              if extension == ".pdf" {
                localfilename :=path.Join(cacheDir, afilename)
                os.MkdirAll(filepath.Dir(localfilename), 0755)
                ioutil.WriteFile(localfilename, a.Content, 0644)
                if smtpconfig.OCRPdf {
                  PutQueue(localfilename+";"+afilename+";"+abucket.Prefix())
                } else {
                  abucket.Upload(localfilename, afilename)
                }
              } else {
                err = abucket.WriteStream(afilename, strings.NewReader(string(a.Content)), 0644)
              }
              if err != nil {
                fmt.Printf("Error writing smtp attachment %s\n", err)
              }
            }
          }
        }
      }
    }()
  }

  return nil
}



func (s *Session) Mail(from string) error {
	fmt.Println("Mail from:", from)
	s.From = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	fmt.Println("Rcpt to:", to)
	emails, err := mail.ParseAddressList(to)
       if err != nil {
           log.Println(err)
           return errors.New("Unable to parse recipients")
       }

       for _, v := range emails {
	       s.To = append(s.To, v.Address)
       }
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func startSMTP() {
  fmt.Println("Starting smtp server 1.3")
  be := &Backend{}

  s := smtp.NewServer(be)

  s.Addr = ":1025"
  s.Domain = "localhost"
  s.ReadTimeout = 300 * time.Second
  s.WriteTimeout = 300 * time.Second
  s.MaxMessageBytes = 90 * 1024 * 1024
  s.MaxRecipients = 50
  s.AllowInsecureAuth = true
        //s.AuthDisabled = true

  cer, err := tls.LoadX509KeyPair("/certs/tls.crt", "/certs/tls.key")
  if err != nil {
    fmt.Println(err)
  } else {
    s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
  }

  if err := s.ListenAndServe(); err != nil {
    fmt.Printf("Error starting smtp %s\n", err)
  }
}


func defineMailFunctions(app *iris.Application, cacheDir string) {
  app.Get("/email/{path:path}", func (c iris.Context) {
    var env *enmime.Envelope
    var err error
    s := mySessions.Start(c)
    fullpath := c.Params().Get("path")
    defaultBucket := getDefaultBucket(s)
    abucket, relativepath := LoadBucketFromPath(fullpath, defaultBucket)
    fmt.Printf("Username: [%s]\n", s.GetString("username"))
    fmt.Printf("Getting full path %s\n", fullpath)
    fmt.Printf("Getting relative path [%s]\n", relativepath)
    reader, err :=  abucket.ReadStream(relativepath)
    if err != nil {
      fmt.Println(err)
    } else {
      env, err = enmime.ReadEnvelope(reader)
      if err != nil {
        fmt.Print(err)
      }
    }
    c.ViewData("Envelope", *env)
    c.ViewData("HeaderKeys", env.GetHeaderKeys())
    header := map[string]string{}
    for _, h := range env.GetHeaderKeys() {
      c.ViewData(h, env.GetHeader(h))
      header[h] = env.GetHeader(h)
    }
    c.ViewData("Headers", header)
    c.ViewData("HTML", env.HTML)
    c.ViewData("Text", env.Text)
    inline := []string{}
    attachments := []string{}
    for _, a := range env.Inlines {
      extension := filepath.Ext(a.FileName)
      afilename := path.Join(filepath.Dir(relativepath), strings.TrimSuffix(a.FileName, extension) + filepath.Base(fullpath) + extension)
      err = abucket.WriteStream(afilename, strings.NewReader(string(a.Content)), 0644)
      if err != nil {
        fmt.Printf("Error writing smtp attachment %s\n", err)
      }
      inline = append(inline, afilename)
    }
    c.ViewData("Inlines", inline)
    for _, a := range env.Attachments {
      extension := filepath.Ext(a.FileName)
      afilename := path.Join(filepath.Dir(relativepath), strings.TrimSuffix(a.FileName, extension) + filepath.Base(fullpath) + extension)
      err = abucket.WriteStream(afilename, strings.NewReader(string(a.Content)), 0644)
      if err != nil {
        fmt.Printf("Error writing smtp attachment %s\n", err)
      }
      attachments = append(attachments, afilename)
    }
    c.ViewData("Attachments", attachments)
    c.ViewData("Username", s.GetString("username"))
    c.ViewData("DefaultBucket", defaultBucket)
    c.ViewData("BucketName", abucket.Name())
    c.ViewData("Bucket", abucket)
    c.ViewData("BucketPrefix", abucket.Prefix())
    c.ViewData("Path", fullpath)
    c.ViewData("RelativePath", relativepath)
    c.ViewData("Dir", filepath.Dir(relativepath))
    c.View("email.html")
    c.Next()
  })
}
