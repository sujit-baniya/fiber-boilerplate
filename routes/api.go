package routes

import (
	Controller "github.com/thomasvvugt/fiber-boilerplate/app/controllers/api"

	"github.com/gofiber/fiber"
)

func RegisterAPI(api *fiber.Group) {
	registerRoles(api)
	registerUsers(api)
}

func registerRoles(api *fiber.Group) {
	roles := api.Group("/roles")

	roles.Get("/", Controller.GetAllRoles)
	roles.Get("/:id", Controller.GetRole)
	roles.Post("/", Controller.AddRole)
	roles.Put("/:id", Controller.EditRole)
	roles.Delete("/:id", Controller.DeleteRole)
}

func registerUsers(api *fiber.Group) {
	users := api.Group("/users")

	users.Get("/", Controller.GetAllUsers)
	users.Get("/:id", Controller.GetUser)
	users.Post("/", Controller.AddUser)
	users.Put("/:id", Controller.EditUser)
	users.Delete("/:id", Controller.DeleteUser)
}
