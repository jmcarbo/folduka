package main

import (
  "testing"
  "folduka/bucket"
)

func TestGetWorkflow(t *testing.T) {
  atestb :=  bucket.NewLocalBucket("test", "./test", "", "", "")
  getWorkflow(atestb, "/", "write", "jmcarbo@imim.es")
}

