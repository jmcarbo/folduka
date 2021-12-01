package main

import (
  "os"
  "github.com/nsqio/go-diskqueue"
  "time"
  "strings"
  "fmt"
)

var dq diskqueue.Interface

func NewTestLogger() diskqueue.AppLogFunc {
    return func(lvl diskqueue.LogLevel, f string, args ...interface{}) {
          fmt.Sprintf(lvl.String()+": "+f, args...)
    }
}

func StartQueue() {
  os.MkdirAll("queueDir", 0755)
  l := NewTestLogger()
  dq = diskqueue.New("queue", "queueDir", 262144, 0, 1<<10, 2500, 2*time.Second, l)
  go func() {
    for {
      select {
      case msgOut := <-dq.ReadChan():
        s := strings.Split(string(msgOut), ";")
        localfilename := s[0]
        afilename := s[1]
        abucketName := s[2]
        fmt.Printf("Loading bucket name %s\n", abucketName)
        abucket := LoadBucket(abucketName)
        fmt.Printf("Loaded bucket %+v\n", abucket)
        if abucket != nil {
          OcrPdf(localfilename)
          err := abucket.Upload(localfilename+".pdf", afilename)
          if err != nil {
            fmt.Printf("Error uploading pdf %s\n", err)
          }
        }
      }
    }
  }()
}

func PutQueue(job string) {
  dq.Put([]byte(job))
}

func StopQueue() {
  dq.Close()
}
