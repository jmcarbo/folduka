package user

import (
  "os"
  "fmt"
  "io/ioutil"
  "path"
  "encoding/json"
  "strings"
)

var userjson = `{ "Username": "%s", "Password": "%s" }`

type LocalUser struct {
  Username, Password string
  DefaultBucket string
}

type LocalUsers struct {
  location string
}

func NewLocalUsers(location string) *LocalUsers {
  os.MkdirAll(location, 0755)
  return &LocalUsers{ location: location }
}

func (u *LocalUsers)Find(username string) *User {
  return nil
}

func (users *LocalUsers)NewUser(username, password string) error {
  os.MkdirAll(path.Join(users.location,username), 0755)
  usrstr := fmt.Sprintf(userjson, username, password)
  ioutil.WriteFile(path.Join(users.location,username+".json"), []byte(usrstr), 0644)
  return nil
}

func (users *LocalUsers)CheckPassword(username, password string) error {
  newusername := strings.ToLower(username)
  b, err := ioutil.ReadFile(path.Join(users.location,newusername+".json"))
  if err != nil {
    return err
  }
  var targetUser LocalUser
  err = json.Unmarshal(b, &targetUser)
  if err != nil {
    return err
  }
  if targetUser.Password != password {
    return fmt.Errorf("Error in user %s validation", newusername)
  }
  return nil
}

func (users *LocalUsers)DefaultBucket(username string) string {
  b, err := ioutil.ReadFile(path.Join(users.location,username+".json"))
  if err != nil {
    return ""
  }
  var targetUser LocalUser
  err = json.Unmarshal(b, &targetUser)
  if err != nil {
    return ""
  }

  return targetUser.DefaultBucket
}
