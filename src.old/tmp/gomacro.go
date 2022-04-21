package main
import (
  "fmt"
  "reflect"
  "github.com/cosmos72/gomacro/fast"
  "github.com/cosmos72/gomacro/imports"
)

func init() {
  imports.Packages["github.com/jordan-wright/email"] = imports.Package{
    Binds:    map[string]reflect.Value{},
    Types:    map[string]reflect.Type{},
    Proxies:  map[string]reflect.Type{},
    Untypeds: map[string]string{},
    Wrappers: map[string][]string{},
  }
}
func RunGomacro(toeval string) reflect.Value {
  interp := fast.New()
  // for simplicity, only collect the first returned value
  val, _ := interp.Eval1(toeval)
  return val
}

const myJson = `
{"data":{"usuaris":[{"cognoms":"jñlkjñlkjñkl","correuElectronic":"","dataInici":"","nom":"jñlkjkl","contrasenya":"","dataInici2":"","usuari":"jñlkjñlkjñkl"}],"textField2":"jkjklkñjkljklj","firstName":"Joe","lastName":"Smith","email":"joe@example.com","submit":true},"metadata":{"timezone":"Europe/Madrid","offset":60,"referrer":"","browserName":"Netscape","userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:64.0) Gecko/20100101 Firefox/64.0","pathName":"/formfill/users.form.5c63116a-2a0b-438e-856b-33df24d3be27.data","onLine":true},"state":"submitted","saved":false}
`

const script = `
import "fmt"
import "github.com/jordan-wright/email"

fmt.Println(myJson)
e := email.NewEmail()
e.From = "Joan Marc Carbo <jmcarbo@gmail.com>"
e.To = []string{"jmcarbo@imim.es"}
//e.Bcc = []string{"test_bcc@example.com"}
//e.Cc = []string{"test_cc@example.com"}
e.Subject = "Awesome Subject"
e.Text = []byte("Text Body is, of course, supported!")
e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "jmcarbo@gmail.com", "2B/nrYsxSBLvS5t5", "smtp.gmail.com"))
`

func main() {
  fmt.Println(RunGomacro("const myJson=`"+myJson+"`\n" + script))
}
