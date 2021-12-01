package bucket

import (
  "github.com/pkg/sftp"
  "golang.org/x/crypto/ssh"
  "fmt"
  "os"
  "io"
  "time"
)


type SFTPElement struct {
  file os.FileInfo
  path string
  bucket Bucket
  prefix string
}

func (e SFTPElement) Prefix() string {
  return e.prefix
}

func (e SFTPElement) Bucket() Bucket {
  return e.bucket
}

func (e SFTPElement) Path() string {
  return e.path
}

func (e SFTPElement) Size() int64 {
  return e.file.Size()
}

func (e SFTPElement) Name() string {
  return e.file.Name()
}

func (e SFTPElement) ModTime() time.Time {
  return e.file.ModTime()
}

func (e SFTPElement) IsDir() bool {
  return e.file.IsDir()
}

func (e SFTPElement) Display() string {
  return "<tr><td>" + e.file.Name() + "</td></tr>"
}

type SFTPBucket struct {
  host string
  user string
  password string
  name string
  root string
  prefix string

  Elements *[]Element

  sftp_destination *sftp.Client
  client_destination *ssh.Client

  Actions map[string]string
  ElementActions map[string]string
}


func NewSFTPBucket(name, root, host, user, password string) *SFTPBucket {
  return &SFTPBucket{ name: name, root: root, host: host, user: user, password: password }
}

func (b SFTPBucket)Name() string {
  return b.name
}

func (b SFTPBucket)Prefix() string {
  return b.prefix
}
func (b *SFTPBucket)SetPrefix(prefix string) {
  b.prefix = prefix
}

func (b SFTPBucket)GetRoot() string {
  return "/"
}

func (b SFTPBucket)DisplayHeader() string {
  return "<table>"
}

func (b SFTPBucket)DisplayFooter() string {
  return "</table>"
}

func (b *SFTPBucket)ListRoot() *[]Element {
  if b.Elements != nil {
    return b.Elements
  }
  els := []Element{}

  err := b.connect()
  defer b.close()
  if err != nil {
    return &els
  }
  files, err := b.sftp_destination.ReadDir(b.root)
  defer b.close()
  if err != nil {
    fmt.Println("'"+b.root+"'")
    return &els
  }
  for _, f := range files {
    els = append(els, SFTPElement{ file: f, bucket: b })
  }
  b.Elements = &els
  return &els
}

func (b *SFTPBucket)List(path string) *[]Element {
  if b.Elements != nil {
    return b.Elements
  }
  els := []Element{}

  err := b.connect()
  defer b.close()
  if err != nil {
    return &els
  }
  files, err := b.sftp_destination.ReadDir(path)
  defer b.close()
  if err != nil {
    fmt.Println("'"+b.root+"'")
    return &els
  }
  for _, f := range files {
    els = append(els, SFTPElement{ file: f, bucket: b })
  }
  b.Elements = &els
  return &els
}

func (b *SFTPBucket)Delete(path string) error {
  return nil
}

func (b *SFTPBucket)Make(path string) error {
  return nil
}

func (b *SFTPBucket)Stat(path string) (os.FileInfo, error) {

  return nil, nil
}

func (b *SFTPBucket)Download(path, local string) error {

  return nil
}

func (b *SFTPBucket)Upload(local, path string) error {

  return nil
}

func (b *SFTPBucket)connect() error {
	var auths []ssh.AuthMethod
        var err error

	auths = append(auths, ssh.Password(b.password))
	config_destination := &ssh.ClientConfig{
		User:            b.user,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	b.client_destination, err = ssh.Dial("tcp", b.host, config_destination)
	if err != nil {
		fmt.Printf("Failed to dial source: " + err.Error())
                return err
	}
	fmt.Println("Successfully connected to ssh server.")

	b.sftp_destination, err = sftp.NewClient(b.client_destination)

	if err != nil {
		fmt.Printf("Error %v\n", err)
                return err
	}

        return nil
}

func (b *SFTPBucket)close() error {
  if b.sftp_destination != nil {
    b.sftp_destination.Close()
  }
  if b.client_destination != nil {
    b.client_destination.Close()
  }
  return nil
}

func (b *SFTPBucket)WriteStream(path string, stream io.Reader, fm os.FileMode) error {
  return nil
}

func (b *SFTPBucket)ReadStream(path string) (io.Reader, error) {
  return nil, nil
}

func (b *SFTPBucket) AddAction(key, value string) {
  b.Actions[key]=value
}

func (b *SFTPBucket) DelAction(key string){
  delete(b.Actions, key)
}

func (b *SFTPBucket) AddElementAction(key, value string){
  b.ElementActions[key]=value
}

func (b *SFTPBucket) DelElementAction(key string){
  delete(b.ElementActions, key)
}

func (b *SFTPBucket)  Rename(source, destination string, isoverwrite bool) error {
  return nil
}

func (b *SFTPBucket)  Copy(source, destination string, isoverwrite bool) error {
  return nil
}

func (b *SFTPBucket)  Search(searchPattern, mypath string) *[]Element {
  return nil
}
