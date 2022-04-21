package user

import (
//  "os"
  "fmt"
  "io/ioutil"
  "path"
  "encoding/json"
  "folduka/bucket"
  "strings"
)


type BucketUser struct {
  Username, Password string
  DefaultBucket string
}

type BucketUsers struct {
  bucket bucket.Bucket
  path string
}

func NewBucketUsers(bucket bucket.Bucket, path string) *BucketUsers {
  return &BucketUsers{ bucket: bucket, path: path }
}

func (u *BucketUsers)Find(username string) *User {
  return nil
}

func (users *BucketUsers)NewUser(username, password string) error {
  return nil
}

func (users *BucketUsers)CheckPassword(username, password string) error {
  newusername := strings.ToLower(username)
  filepath := path.Join(users.path, "_config", newusername+".json")
  fmt.Printf("Reading file %s\n", filepath)
  reader, err := users.bucket.ReadStream(filepath)
  if err != nil {
    fmt.Printf("User %s not found in path %s with error %s\n", username, path.Join(users.path, "_config", newusername+".json"), err)
    return err
  }
  b, err := ioutil.ReadAll(reader)
  if err != nil {
    return err
  }
  var targetUser BucketUser
  err = json.Unmarshal(b, &targetUser)
  if err != nil {
    return err
  }
  if targetUser.Password != password {
    return fmt.Errorf("Error in user %s validation", username)
  }
  fmt.Printf("user [%s] validated OK\n")
  return nil
}

func (users *BucketUsers)DefaultBucket(username string) string {
  reader, err := users.bucket.ReadStream(path.Join(users.path, "_config", username+".json"))
  if err != nil {
    return ""
  }
  b, err := ioutil.ReadAll(reader)
  if err != nil {
    return ""
  }
  var targetUser BucketUser
  err = json.Unmarshal(b, &targetUser)
  if err != nil {
    return "" 
  }

  return targetUser.DefaultBucket
}
