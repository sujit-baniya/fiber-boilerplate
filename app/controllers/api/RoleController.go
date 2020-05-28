package api

import (
	"github.com/gofiber/fiber"

	"github.com/thomasvvugt/fiber-boilerplate/app/models"
	"github.com/thomasvvugt/fiber-boilerplate/database"
)

// Return all roles as JSON
func GetAllRoles(c *fiber.Ctx) {
	db := database.Instance()
	var Role []models.Role
	if res := db.Find(&Role); res.Error != nil {
		c.Send("Error occurred while retrieving roles from the database", res.Error)
	}
	err := c.JSON(Role)
	if err != nil {
		panic("Error occurred when returning JSON of roles")
	}
}

// Return a single role as JSON
func GetRole(c *fiber.Ctx) {
	db := database.Instance()
	Role := new(models.Role)
	id := c.Params("id")
	if res := db.Find(&Role, id); res.Error != nil {
		c.Send("An error occurred when retrieving the role", res.Error)
	}
	if Role.ID == 0 {
		c.SendStatus(404)
		err := c.JSON(fiber.Map{
			"ID": id,
		})
		if err != nil {
			panic("Error occurred when returning JSON of a role")
		}
		return
	}
	err := c.JSON(Role)
	if err != nil {
		panic("Error occurred when returning JSON of a role")
	}
}

// Add a single role to the database
func AddRole(c *fiber.Ctx) {
	db := database.Instance()
	Role := new(models.Role)
	if err := c.BodyParser(Role); err != nil {
		c.Send("An error occurred when parsing the new role", err)
	}
	if res := db.Create(&Role); res.Error != nil {
		c.Send("An error occurred when storing the new role", res.Error)
	}
	err := c.JSON(Role)
	if err != nil {
		panic("Error occurred when returning JSON of a role")
	}
}

// Edit a single role
func EditRole(c *fiber.Ctx) {
	db := database.Instance()
	id := c.Params("id")
	EditRole := new(models.Role)
	Role := new(models.Role)
	if err := c.BodyParser(EditRole); err != nil {
		c.Send("An error occurred when parsing the edited role", err)
	}
	if res := db.Find(&Role, id); res.Error != nil {
		c.Send("An error occurred when retrieving the existing role", res.Error)
	}
	// Role does not exist
	if Role.ID == 0 {
		c.SendStatus(404)
		err := c.JSON(fiber.Map{
			"ID": id,
		})
		if err != nil {
			panic("Error occurred when returning JSON of a role")
		}
		return
	}
	Role.Name = EditRole.Name
	Role.Description = EditRole.Description
	db.Save(&Role)

	err := c.JSON(Role)
	if err != nil {
		panic("Error occurred when returning JSON of a role")
	}
}

// Delete a single role
func DeleteRole(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.Instance()

	var Role models.Role
	db.Find(&Role, id)
	if res := db.Find(&Role); res.Error != nil {
		c.Send("An error occurred when finding the role to be deleted", res.Error)
	}
	db.Delete(&Role)

	err := c.JSON(fiber.Map{
		"ID": id,
		"Deleted": true,
	})
	if err != nil {
		panic("Error occurred when returning JSON of a role")
	}
}
