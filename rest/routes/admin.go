package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/fiber-boilerplate/rest/controllers"
)

func AdminRoutes(web fiber.Router) {
	admin := web.Group("admin")
	admin.Get("/", controllers.Admin)
	admin.Get("/users", controllers.UserList)
	AdminUserRoutes(admin)
	RolesRoutes(admin)
	PermissionRoutes(admin)
}

func AdminUserRoutes(a fiber.Router) {
	services := a.Group("users")
	services.Get("/", controllers.UserList)
	services.Get("/:id", controllers.UserInfo)
	services.Put("/:id", controllers.UpdateUser)
	services.Get("/:id/settings", controllers.UserSettings)
	services.Post("/:id/settings", controllers.StoreUserSettings)
}

func RolesRoutes(r fiber.Router) {
	roles := r.Group("roles")
	roles.Post("/create", controllers.CreateNewRole)
	roles.Post("/remove", controllers.RemoveRole)
	roles.Post("/assign", controllers.AssignRoleToUser)
	roles.Post("/revoke", controllers.RevokeRoleFromUser)
	roles.Post("/change", controllers.ChangeRoleForUser)
}

func PermissionRoutes(r fiber.Router) {
	roles := r.Group("permissions")
	roles.Post("/add", controllers.AddPermissionOnRole)
	roles.Post("/remove", controllers.RemovePermissionFromRole)
	roles.Post("/change", controllers.ChangePermissionOnRole)
}
