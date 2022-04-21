package user

type User interface {
}

type Users interface {
  Find(username string) *User
  NewUser(username, password string) error
  CheckPassword(password string) error
  DefaultBucket(username string) string
}
