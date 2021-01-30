package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/fiber-boilerplate/app"
	"github.com/sujit-baniya/fiber-boilerplate/rest/controllers"
	"github.com/sujit-baniya/fiber-boilerplate/rest/middlewares"
)

func UserRoutes(web fiber.Router) {

	web.Use(middlewares.AuthWeb())
	if app.Http.Auth.Type == "casbin" {
		web.Use(app.Http.Auth.Casbin.RoutePermission())
	}
	account := web.Group("/app")
	account.Get("/", controllers.App)
	account.Get("/me", controllers.Me)
	FileRoutes(account)
	PaymentRoutes(account)
	AdminRoutes(account)
}

func FileRoutes(account fiber.Router) {
	account.Get("/files", controllers.FileIndex)
	account.Post("/do/upload", middlewares.MaxBodySize(20), controllers.Upload)
}

func PaymentRoutes(account fiber.Router) {
	account.Post("/paypal/do/order", controllers.PlaceOrderFromPaypal)
	account.Post("/paypal/do/order/validate/:amount", controllers.ValidateOrderFromPaypal)
	account.Post("/paypal/order/success/:id", controllers.PostOrderSuccessResponseFromPaypal)
	account.Get("/paypal/order/success/:id", controllers.PostOrderSuccessResponseFromPaypal)
	account.Post("/paypal/order/cancel/:id", controllers.PostOrderCancelResponseFromPaypal)
	account.Get("/paypal/order/:id", controllers.GetOrderDetailFromPaypal)
}
