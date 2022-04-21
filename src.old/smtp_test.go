package main

import (
  "testing"
  "log"
  "io/ioutil"
  "encoding/json"
)

func TestReadSMTPFolderConfig(t *testing.T) {
  b, err:=ioutil.ReadFile("files/jmca/_config/smtp.json")
  if err !=nil {
    log.Fatal(err)
  }
  c := smtpFolderConfig{}
  err = json.Unmarshal(b, &c)
  if err !=nil {
    log.Fatal(err)
  }
}
