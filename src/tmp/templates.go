import (
  //"fmt"
  "text/template"
  "encoding/json"
  "os"
)


func main() {
  t := template.Must(template.New("").Parse(templ))

  m := map[string]interface{}{}
  if err := json.Unmarshal([]byte(jsondata), &m); err != nil {
    panic(err)
  }

  if err := t.Execute(os.Stdout, m); err != nil {
    panic(err)
  }
}

const templ = `<html><body>
Value of a: {{ range .data.usuaris}}{{.cognoms}}{{end}}
Something else: {{.metadata}}
</body></html>`

const jsondata = `
{"data":{"usuaris":[{"cognoms":"jñlkjñlkjñkl","correuElectronic":"","dataInici":"","nom":"jñlkjkl","contrasenya":"","dataInici2":"","usuari":"jñlkjñlkjñkl"}],"textField2":"jkjklkñjkljklj","firstName":"Joe","lastName":"Smith","email":"joe@example.com","submit":true},"metadata":{"timezone":"Europe/Madrid","offset":60,"referrer":"","browserName":"Netscape","userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.14; rv:64.0) Gecko/20100101 Firefox/64.0","pathName":"/formfill/users.form.5c63116a-2a0b-438e-856b-33df24d3be27.data","onLine":true},"state":"submitted","saved":false}

`
