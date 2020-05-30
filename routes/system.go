package routes

import (
	"github.com/gofiber/fiber"
	"github.com/thomasvvugt/fiber-boilerplate/routes/system"
)

func RegisterSystemRoutes(app *fiber.Group) {
	app.Get("/", system.SystemInfo)
	app.Get("/cpu", system.CpuInfo)
	app.Get("/disk", system.DiskInfo)
	app.Get("/memory", system.RamInfo)
	app.Get("/host", system.HostInfo)
	app.Get("/processes", system.ProcessInfo)
	app.Get("/internet", system.InternetInfo)
	app.Get("/bandwidth", system.BandwidthInfo)
	app.Get("/db-status", system.CheckDB)
}
