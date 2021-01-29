package controllers

import (
	"github.com/gofiber/fiber/v2" //nolint:goimports
	"github.com/sujit-baniya/fiber-boilerplate/pkg/auth"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/models"
	"strconv"
)

func UserList(c *fiber.Ctx) error {
	users := models.AllUsers()
	return c.JSON(users)
}
func UserInfo(c *fiber.Ctx) error {
	users, err := models.GetUserById(c.Params("id"))
	if err != nil {
		return c.JSON(fiber.Map{
			"message": err.Error(),
			"error":   true,
		})
	}
	return c.JSON(users)
}
func UpdateUser(c *fiber.Ctx) error {
	var u models.User
	uid, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	c.BodyParser(&u)
	u.ID = uint(uid)
	u.Update()
	return c.JSON(u)
}

func Me(c *fiber.Ctx) error {
	user, _ := auth.User(c)
	return c.JSON(user)
}

func UserSettings(c *fiber.Ctx) error {
	uid, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	user := models.User{
		ID: uint(uid),
	}
	settings, err := user.Settings()
	if err != nil {
		return c.JSON(fiber.Map{
			"message": err.Error(),
			"error":   true,
		})
	}
	return c.JSON(settings)
}

func StoreUserSettings(c *fiber.Ctx) error {
	var userSettings models.UserSetting
	c.BodyParser(&userSettings)
	uid, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	userSettings.UserID = uint(uid)
	userSettings.UpdateOrCreate()
	return c.JSON(userSettings)
}
