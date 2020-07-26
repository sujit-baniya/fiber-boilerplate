package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber"
	. "github.com/itsursujit/fiber-boilerplate/app"
	"github.com/itsursujit/fiber-boilerplate/config"
	"github.com/itsursujit/fiber-boilerplate/libraries"
	"github.com/itsursujit/fiber-boilerplate/middlewares"
	"github.com/itsursujit/fiber-boilerplate/models"
	"github.com/itsursujit/fiber-boilerplate/routes"
)

func main() {
	Log = libraries.SetupZeroLog()
	migrate := flag.Bool("migrate", false, "Migrate the pending migration files")
	flag.Parse()
	if *migrate {
		initMigrate()
		return
	}
	Serve()
}

func Serve() {
	Boot()
	App.Use(middlewares.LogMiddleware)
	routes.Load()
	App.Use(func(c *fiber.Ctx) {
		var err fiber.Error
		err.Code = fiber.StatusNotFound
		config.CustomErrorHandler(c, &err)
	})
	// go libraries.Consume("webhook-callback")               //nolint:wsl
	err := App.Listen(config.AppConfig.App_Port) //nolint:wsl
	if err != nil {
		panic("App not starting: " + err.Error() + "on Port: " + config.AppConfig.App_Port)
	}

	defer DB.Close()
}

func Boot() {
	config.LoadEnv()
	config.BootApp()
	LoadComponents()
}

func LoadComponents() {
	config.LoadQueueConfig()
	config.LoadPaypalConfig()
	Queue = libraries.SetupQueue() //nolint:wsl
}

func initMigrate() {
	fmt.Println("1")
	config.LoadEnv()
	_, err := config.SetupDB()
	if err != nil {
		panic(err)
	}
	Migrate()
}

func Migrate() {
	fmt.Println("Migrating...")
	Log.Info().Msg("Migrating")
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.PaymentMethod{})
	DB.AutoMigrate(&models.Payment{})
	DB.AutoMigrate(&models.Transaction{})
	DB.AutoMigrate(&models.UserTransactionLog{})
	DB.AutoMigrate(&models.File{})
	DB.AutoMigrate(&models.UserFile{})
	Log.Info().Msg("Migrated")
	fmt.Println("Migrated...")
}
