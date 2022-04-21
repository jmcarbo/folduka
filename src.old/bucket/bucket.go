// Package bucket
package bucket

import (
  "io"
  "os"
  "strings"
  "time"
)

type Element interface {
  Name() string
  Path() string
  Prefix() string
  Display() string
  IsDir() bool
  ModTime() time.Time
  Size() int64
  Bucket() Bucket
}

type ElementByName []Element

func (e ElementByName) Len() int {
  return len(e)
}

func (e ElementByName) Swap(i, j int) {
  e[i], e[j] = e[j], e[i]
}

func (e ElementByName) Less(i, j int) bool {
  if e[i].Name() == "_config" {
    return true
  }
  if strings.Compare(strings.ToLower(e[i].Name()), strings.ToLower(e[j].Name())) < 0 {
    return true
  }
  return false
}

type Bucket interface {
  Name() string
  SetPrefix(prefix string)
  Prefix() string
  GetRoot() string
  Make(path string) error
  Delete(path string) error
  Rename(source, destination string, isoverwrite bool) error
  Copy(source, destination string, isoverwrite bool) error
  ListRoot() *[]Element
  List(path string) *[]Element
  DisplayHeader() string
  DisplayFooter() string
  Stat(path string) (os.FileInfo, error)
  WriteStream(path string, stream io.Reader, _ os.FileMode) error
  ReadStream(path string) (io.Reader, error)
  Download(path, local string) error
  Upload(local, path string) error
  AddAction(key, value string)
  DelAction(key string)
  AddElementAction(key, value string)
  DelElementAction(key string)
  Search(searchPattern, mypath string) *[]Element
}

