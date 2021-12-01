package main

import (
  "folduka/bucket"
  "bytes"
  "encoding/json"
  "path"
  "fmt"
)



type Display struct {
  Mask string
  Template string
}

func getDisplay(abucket bucket.Bucket, mypath string) ([]Display, error) {
  fmt.Printf("Getting display from bucket %s with path %s\n", abucket.Name(), mypath)
  r, err := abucket.ReadStream(path.Join(mypath, "_config/display.json"))
  if err != nil {
    r, err = abucket.ReadStream("_config/display.json")
    if err != nil {
      return nil, err
    }
  }

  buf := new(bytes.Buffer)
  buf.ReadFrom(r)
  display := []Display{}
  err = json.Unmarshal(buf.Bytes(), &display)
  if err != nil {
    return nil, err
  }
  return display, nil
}
