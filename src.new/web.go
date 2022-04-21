package main

import (
	"time"
	//"github.com/dchest/captcha"
	"github.com/jmcarbo/folduka/controllers"
	"github.com/rs/zerolog/log"
	//"net/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"html/template"
	"github.com/jmcarbo/form"
	"strings"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var address = ":8080"

var inputTpl = `
<div class="mb-4">
	<label class="block text-grey-darker text-sm font-bold mb-2" {{with .ID}}for="{{.}}"{{end}}>
		{{.Label}}
	</label>
	<input class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight" {{with .ID}}id="{{.}}"{{end}} type="{{.Type}}" name="{{lower .Name}}" id="{{ lower .Name }}" placeholder="{{.Placeholder}}" {{with .Value}}value="{{.}}"{{end}}>
	{{with .Footer}}
		<p class="text-grey pt-2 text-xs italic">{{.}}</p>
	{{end}}
</div>`

func initWebserver() {

	tpl := template.Must(template.New("").Funcs(template.FuncMap{
		"lower": strings.ToLower}).Parse(inputTpl))
	fb := form.Builder{
		InputTemplate: tpl,
	}

	log.Info().Msg("Starting server at " + address)

	engine := html.New("./views", ".html")
	engine.AddFunc("inputs_for", fb.Inputs)
	engine.AddFunc("inputs_and_errors_for", func(v interface{}, errs []error) (template.HTML, error) {
			return fb.Inputs(v, errs...)
		})
	engine.AddFunc("lower", strings.ToLower)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	
	controllers.Store = session.New()

	// Initialize default config
	app.Use(csrf.New(csrf.Config{
    		//KeyLookup:      "header:X-Csrf-Token",
    		KeyLookup:      "form:csrf",
    		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		ContextKey:     "token",
	}))

	app.Static("/js", "public/js")
	app.Static("/css", "public/css")
	//app.Get("/", controllers.InitDashboard)
	//app.Post("/", controllers.Dashboard)

	app.Get("/login", controllers.InitLogin)
	app.Post("/login", controllers.Login)

	app.Listen(address)
}
