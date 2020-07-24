package middlewares

import (
	"github.com/gofiber/fiber"
	"github.com/gookit/validate"
	. "github.com/itsursujit/fiber-boilerplate/app"
	"github.com/itsursujit/fiber-boilerplate/config"
	"github.com/itsursujit/fiber-boilerplate/libraries"
	"github.com/itsursujit/fiber-boilerplate/models"
)

func ValidateRegisterPost(c *fiber.Ctx) {
	var register models.RegisterForm
	if err := c.BodyParser(&register); err != nil {
		Flash.WithError(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/register")
		return
	}

	v := validate.Struct(register)
	if !v.Validate() {
		Flash.WithError(c, fiber.Map{
			"message": v.Errors.One(),
		}).Redirect("/register")
		return
	}
	c.Locals("register", register)
	c.Next()
}

func ValidateConfirmToken(c *fiber.Ctx) {
	t := libraries.Decrypt(c.Query("t"), config.AppConfig.App_Key)
	user, err := models.GetUserByEmail(t)
	if err != nil {
		Flash.WithError(c, fiber.Map{
			"message": err.Error(),
		}).Redirect("/login")
		return
	}

	if user.EmailVerified {
		Flash.WithError(c, fiber.Map{
			"message": "Email was already validated",
		}).Redirect("/login")
		return
	}
	c.Locals("user", user)
	c.Next()
}
