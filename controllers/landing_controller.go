package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/fiber-boilerplate/auth"
)

func Landing(c *fiber.Ctx) error {
	user, _ := auth.User(c)
	layout := "layouts/main"
	view := "index"
	if user == nil {
		layout = "layouts/landing"
		view = "landing"
	}

	return c.Render(view, fiber.Map{
		"auth": user != nil,
		"user": user,
	}, layout)
}
