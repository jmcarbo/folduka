package main

import (
  "testing"
  "folduka/bucket"
)

func TestSearch(t *testing.T) {
  abucket :=  bucket.NewLocalBucket("local","./files",
      "", "", "" )
  s:=abucket.Search("jmcarbo", "")
  for _, f := range *s {
    if f != nil {
      t.Log(f.Name())
      t.Log(f.Path())
    }
  }
}
