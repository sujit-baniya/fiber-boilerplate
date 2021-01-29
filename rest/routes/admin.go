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
}

func AdminUserRoutes(a fiber.Router) {
	services := a.Group("users")
	services.Get("/", controllers.UserList)
	services.Get("/:id", controllers.UserInfo)
	services.Put("/:id", controllers.UpdateUser)
	services.Get("/:id/settings", controllers.UserSettings)
	services.Post("/:id/settings", controllers.StoreUserSettings)
}

func CasbinRoutes(a fiber.Router) {

}
