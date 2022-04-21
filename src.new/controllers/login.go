package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmcarbo/folduka/models"
	"log"
	"strings"
)


func InitLogin(c *fiber.Ctx) error {
	return c.Render("login", 
		fiber.Map{ "login": models.LoginForm{}, 
			"csrf": c.Locals("token") }, 
			"layouts/main")
}


func Login(c *fiber.Ctx) error {
        p := new(models.LoginForm)
        if err := c.BodyParser(p); err != nil {
            return err
        }

	log.Printf(">>> %+v\n", p)
	if models.ValidateByEmail(p.Email, p.Password) {
		sess, _ := Store.Get(c)
		sess.Set("email", strings.ToLower(p.Email))
		err := sess.Save()
		if err != nil {
			log.Printf("Error saving session: %s", err.Error())
		}
		return c.Redirect("/")	
	} else {
		return c.Render("login", 
			fiber.Map{ "login": models.LoginForm{}, 
			"csrf": c.Locals("token") }, 
			"layouts/main")
	}
}

