package controllers

import (
	"github.com/gofiber/fiber/v2" //nolint:goimports
	. "github.com/itsursujit/fiber-boilerplate/app"
	"github.com/itsursujit/fiber-boilerplate/auth"
	"github.com/itsursujit/fiber-boilerplate/config"
	"github.com/itsursujit/fiber-boilerplate/models"
)

func LoginGet(c *fiber.Ctx) error {
	Flash.Get(c)
	return c.Render("auth/login", Flash.Data, "layouts/auth")
}

func LoginPost(c *fiber.Ctx) error { //nolint:wsl
	user := c.Locals("user").(*models.User)
	_, _ = auth.Login(c, user.ID, config.AuthConfig.App_Jwt_Secret) //nolint:wsl
	return c.Redirect("/")
}

func LogoutPost(c *fiber.Ctx) error { //nolint:nolintlint,wsl
	if auth.IsLoggedIn(c) {
		_ = auth.Logout(c)
	}
	return c.Redirect("/login")
}
