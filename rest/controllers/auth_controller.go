package controllers

import (
	"github.com/gofiber/fiber/v2" //nolint:goimports
	"github.com/sujit-baniya/fiber-boilerplate/app"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/auth"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/models"
)

func LoginGet(c *fiber.Ctx) error {
	data := getFlashData(c)
	data["title"] = "Login | "
	return c.Render("auth/login", data, "layouts/landing")
}

func LoginPost(c *fiber.Ctx) error { //nolint:wsl
	user := c.Locals("user").(*models.User)
	_, err := auth.Login(c, user.ID, app.Http.Token.AppJwtSecret) //nolint:wsl
	if err != nil {
		return app.Http.Flash.WithError(c, fiber.Map{
			"error":   true,
			"message": err.Error(),
		}).Redirect("/login")
	}
	return c.Redirect("/")
}

func LogoutPost(c *fiber.Ctx) error { //nolint:nolintlint,wsl
	if auth.IsLoggedIn(c) {
		err := auth.Logout(c)
		if err != nil {
			panic(err)
		}
	}
	c.Set("X-DNS-Prefetch-Control", "off")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
	c.Set("Cache-Control", "no-cache, must-revalidate, no-store, max-age=0, private")
	return c.Redirect("/login")
}

func getFlashData(c *fiber.Ctx) fiber.Map {
	data := app.Http.Flash.Get(c)
	if data == nil {
		data = fiber.Map{}
	}
	return data
}
