package routes

import (
	"github.com/gofiber/fiber/v2"
	. "github.com/sujit-baniya/fiber-boilerplate/app"
	"github.com/sujit-baniya/fiber-boilerplate/auth"
	"github.com/sujit-baniya/fiber-boilerplate/config"
	"github.com/sujit-baniya/fiber-boilerplate/controllers"
	"github.com/sujit-baniya/fiber-boilerplate/middlewares"
)

func WebRoutes() {
	web := App.Group("")
	web.Use(auth.AuthCookie)
	LandingRoutes(web)
	UserRoutes(web)
}

func LandingRoutes(app fiber.Router) {
	app.Get("/", controllers.Landing)
}

func UserRoutes(app fiber.Router) {
	account := app.Group("/account")
	account.Use(middlewares.Authenticate(middlewares.AuthConfig{
		SigningKey:  []byte(config.AuthConfig.App_Jwt_Secret),
		TokenLookup: "cookie:fiber-boilerplate-Token",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			auth.Logout(ctx)
			return ctx.Redirect("/login")
		},
	}))
	// account.Get("/users", controllers.Index)
	account.Get("/file-manager", controllers.FileIndex)
	account.Get("/file-manager/view", controllers.ViewFile)
	account.Post("/file-manager/upload", controllers.Upload)
	account.Post("/paypal/do/order", controllers.PlaceOrderFromPaypal)
	account.Post("/paypal/do/order/validate/:amount", controllers.ValidateOrderFromPaypal)
	account.Get("/paypal/order/success/:id", controllers.PostOrderSuccessResponseFromPaypal)
	account.Post("/paypal/order/cancel/:id", controllers.PostOrderCancelResponseFromPaypal)
	account.Get("/paypal/order/:id", controllers.GetOrderDetailFromPaypal)
}
