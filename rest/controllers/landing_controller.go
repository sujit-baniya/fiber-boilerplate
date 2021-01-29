package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/auth"
)

func Landing(c *fiber.Ctx) error {
	user, err := auth.User(c)
	if err != nil {
		auth.Logout(c)
	}
	view := "index"

	if err := c.Render(view, fiber.Map{
		"auth": user != nil,
		"user": user,
	}, "layouts/landing"); err != nil {
		panic(err.Error())
	}
	return nil
}

func Terms(c *fiber.Ctx) error {
	if err := c.Render("terms", fiber.Map{"title": "Terms | "}, "layouts/landing"); err != nil {
		panic(err.Error())
	}
	return nil
}

func PrivacyPolicy(c *fiber.Ctx) error {
	if err := c.Render("privacy-policy", fiber.Map{"title": "Privacy Policy | "}, "layouts/landing"); err != nil {
		panic(err.Error())
	}
	return nil
}

func Disclaimer(c *fiber.Ctx) error {
	if err := c.Render("disclaimer", fiber.Map{"title": "Disclaimer | "}, "layouts/landing"); err != nil {
		panic(err.Error())
	}
	return nil
}

func App(c *fiber.Ctx) error {
	user, err := auth.User(c)
	if err != nil {
		auth.Logout(c)
	}
	view := "home"

	if err := c.Render(view, fiber.Map{
		"auth": user != nil,
		"user": user,
	}); err != nil {
		panic(err.Error())
	}
	return nil
}

func Admin(c *fiber.Ctx) error {
	user, err := auth.User(c)
	if err != nil {
		auth.Logout(c)
	}
	view := "admin-home"

	if err := c.Render(view, fiber.Map{
		"auth": user != nil,
		"user": user,
	}); err != nil {
		panic(err.Error())
	}
	return nil
}
