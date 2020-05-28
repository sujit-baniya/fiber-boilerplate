package api

import (
	"github.com/gofiber/fiber"

	"github.com/thomasvvugt/fiber-boilerplate/app/models"
	"github.com/thomasvvugt/fiber-boilerplate/database"
)

// Return all users as JSON
func GetAllUsers(c *fiber.Ctx) {
	db := database.Instance()
	var Users []models.User
	if res := db.Find(&Users); res.Error != nil {
		c.Send("Error occurred while retrieving users from the database", res.Error)
	}
	// Match roles to users
	for index, User := range Users {
		if User.RoleID != 0 {
			Role := new(models.Role)
			if res := db.Find(&Role, User.RoleID); res.Error != nil {
				c.Send("An error occurred when retrieving the role", res.Error)
			}
			if Role.ID != 0 {
				Users[index].Role = *Role
			}
		}
	}
	err := c.JSON(Users)
	if err != nil {
		panic("Error occurred when returning JSON of users")
	}
}

// Return a single user as JSON
func GetUser(c *fiber.Ctx) {
	db := database.Instance()
	User := new(models.User)
	id := c.Params("id")
	if res := db.Find(&User, id); res.Error != nil {
		c.Send("An error occurred when retrieving the user", res.Error)
	}
	if User.ID == 0 {
		c.SendStatus(404)
		err := c.JSON(fiber.Map{
			"ID": id,
		})
		if err != nil {
			panic("Error occurred when returning JSON of a role")
		}
		return
	}
	// Match role to user
	if User.RoleID != 0 {
		Role := new(models.Role)
		if res := db.Find(&Role, User.RoleID); res.Error != nil {
			c.Send("An error occurred when retrieving the role", res.Error)
		}
		if Role.ID != 0 {
			User.Role = *Role
		}
	}
	err := c.JSON(User)
	if err != nil {
		panic("Error occurred when returning JSON of a user")
	}
}

// Add a single user to the database
func AddUser(c *fiber.Ctx) {
	db := database.Instance()
	User := new(models.User)
	if err := c.BodyParser(User); err != nil {
		c.Send("An error occurred when parsing the new user", err)
	}
	if res := db.Create(&User); res.Error != nil {
		c.Send("An error occurred when storing the new user", res.Error)
	}
	// Match role to user
	if User.RoleID != 0 {
		Role := new(models.Role)
		if res := db.Find(&Role, User.RoleID); res.Error != nil {
			c.Send("An error occurred when retrieving the role", res.Error)
		}
		if Role.ID != 0 {
			User.Role = *Role
		}
	}
	err := c.JSON(User)
	if err != nil {
		panic("Error occurred when returning JSON of a user")
	}
}

// Edit a single user
func EditUser(c *fiber.Ctx) {
	db := database.Instance()
	id := c.Params("id")
	EditUser := new(models.User)
	User := new(models.User)
	if err := c.BodyParser(EditUser); err != nil {
		c.Send("An error occurred when parsing the edited user", err)
	}
	if res := db.Find(&User, id); res.Error != nil {
		c.Send("An error occurred when retrieving the existing user", res.Error)
	}
	// User does not exist
	if User.ID == 0 {
		c.SendStatus(404)
		err := c.JSON(fiber.Map{
			"ID": id,
		})
		if err != nil {
			panic("Error occurred when returning JSON of a user")
		}
		return
	}
	User.Name = EditUser.Name
	User.Email = EditUser.Email
	User.RoleID = EditUser.RoleID
	// Match role to user
	if User.RoleID != 0 {
		Role := new(models.Role)
		if res := db.Find(&Role, User.RoleID); res.Error != nil {
			c.Send("An error occurred when retrieving the role", res.Error)
		}
		if Role.ID != 0 {
			User.Role = *Role
		}
	}
	// Save user
	db.Save(&User)

	err := c.JSON(User)
	if err != nil {
		panic("Error occurred when returning JSON of a user")
	}
}

// Delete a single user
func DeleteUser(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.Instance()

	var User models.User
	db.Find(&User, id)
	if res := db.Find(&User); res.Error != nil {
		c.Send("An error occurred when finding the user to be deleted", res.Error)
	}
	db.Delete(&User)

	err := c.JSON(fiber.Map{
		"ID": id,
		"Deleted": true,
	})
	if err != nil {
		panic("Error occurred when returning JSON of a user")
	}
}
