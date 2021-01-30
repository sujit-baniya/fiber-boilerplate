package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/sujit-baniya/fiber-boilerplate/app"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/models"
	"gorm.io/gorm"
)

func CreateNewRole(c *fiber.Ctx) error {
	var role models.Role
	c.BodyParser(&role)
	role.Slug = slug.Make(role.Name)
	app.Http.Database.Create(&role)
	return c.JSON(role)
}
func RemoveRole(c *fiber.Ctx) error {
	var role models.Role
	c.BodyParser(&role)
	role.Slug = slug.Make(role.Name)
	app.Http.Database.
		Delete(&models.RoleAndPermission{}).
		Where(models.RoleAndPermission{V0: role.Slug}).
		Or(models.RoleAndPermission{V1: role.Slug})
	app.Http.Auth.Enforcer.LoadPolicy()
	return nil
}

func AssignRoleToUser(c *fiber.Ctx) error {
	var roleRequest models.RoleRequest
	var role1 models.RoleAndPermission
	c.BodyParser(&roleRequest)
	role := models.RoleAndPermission{
		Ptype: "g",
		V0: fmt.Sprintf("%d", roleRequest.UserID),
		V1: roleRequest.Role,
	}
	err := app.Http.Database.Unscoped().First(&role1, models.RoleAndPermission{Ptype: role.Ptype,V0: role.V0, V1: role.V1}).Error
	if err != nil {
		app.Http.Database.Create(&role)
		app.Http.Auth.Enforcer.LoadPolicy()
		return c.JSON(role)
	}
	role1.DeletedAt = gorm.DeletedAt{Valid: false}
	app.Http.Database.Unscoped().Save(role1)
	app.Http.Auth.Enforcer.LoadPolicy()
	return c.JSON(role1)
}

func RevokeRoleFromUser(c *fiber.Ctx) error {
	var roleRequest models.RoleRequest
	c.BodyParser(&roleRequest)
	app.Http.Database.
		Delete(
			&models.RoleAndPermission{},
			models.RoleAndPermission{V0: fmt.Sprintf("%d", roleRequest.UserID), V1: roleRequest.Role})
	app.Http.Auth.Enforcer.LoadPolicy()
	return c.JSON("Role Revoked from user")
}

func ChangeRoleForUser(c *fiber.Ctx) error {
	var role models.RoleRequest
	var role1 models.RoleAndPermission
	c.BodyParser(&role)
	err := app.Http.Database.Unscoped().First(&role1, models.RoleAndPermission{Ptype: "g",V0: fmt.Sprintf("%d", role.UserID), V1: role.OldRole}).Error
	if err == nil {
		role1.V1 = role.Role
		role1.DeletedAt = gorm.DeletedAt{Valid: false}
		app.Http.Database.Unscoped().Save(role1)
		app.Http.Auth.Enforcer.LoadPolicy()
		return c.JSON("Role Changed for user")
	}
	return c.JSON("Role doesn't exists")
}

func AddPermissionOnRole(c *fiber.Ctx) error {
	var permission models.PermissionRequest
	var role1 models.RoleAndPermission
	c.BodyParser(&permission)
	if permission.Role != "" && permission.Module != "" && permission.Action != "" {
		role := models.RoleAndPermission{
			Ptype: "p",
			V0: permission.Role,
			V1: permission.Module,
			V2: permission.Action,
			Category: "permission",
		}
		err := app.Http.Database.Unscoped().First(&role1, models.RoleAndPermission{Ptype: role.Ptype,V0: role.V0, V1: role.V1, V2: role.V2}).Error
		if err != nil {
			app.Http.Database.Create(&role)
		}
		role1.DeletedAt = gorm.DeletedAt{Valid: false}
		app.Http.Database.Unscoped().Save(role1)
	}
	if permission.Role != "" && permission.Route != "" && permission.Method != "" {
		fmt.Println(1)
		role := models.RoleAndPermission{
			Ptype: "p",
			V0: permission.Role,
			V1: permission.Route,
			V2: permission.Method,
			Category: "route",
		}
		err := app.Http.Database.Unscoped().First(&role1, models.RoleAndPermission{Ptype: role.Ptype,V0: role.V0, V1: role.V1, V2: role.V2}).Error
		if err != nil {
			app.Http.Database.Create(&role)
		}
		role1.DeletedAt = gorm.DeletedAt{Valid: false}
		app.Http.Database.Unscoped().Save(role1)
	}
	app.Http.Auth.Enforcer.LoadPolicy()
	return nil
}
func RemovePermissionFromRole(c *fiber.Ctx) error {
	return nil
}
func ChangePermissionOnRole(c *fiber.Ctx) error {
	return nil
}
