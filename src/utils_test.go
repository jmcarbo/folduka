package main

import (
  "testing"
  "fmt"
)

func TestOCR(t *testing.T) {
  err := OcrPdf("cache/sg/CopiaAutentica/bla.pdf")
  if err != nil {
    fmt.Println(err)
  }
}


