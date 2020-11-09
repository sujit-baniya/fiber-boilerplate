package controllers

import (
	"github.com/gofiber/fiber/v2" //nolint:goimports
	. "github.com/sujit-baniya/fiber-boilerplate/app"
	"github.com/sujit-baniya/fiber-boilerplate/models"
)

func Index(c *fiber.Ctx) error {
	var users []models.User
	DB.Find(&users)   //nolint:wsl
	return c.JSON(users) //nolint:errcheck
}
