package system

import (
	"fmt"
	"runtime"

	"github.com/gofiber/fiber"
	"github.com/shirou/gopsutil/host"
)

type Host struct {
	Name     string `json:"name"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Platform string `json:"platform"`
	Uptime   uint64 `json:"uptime"`
}

func CheckHost() Host {
	host, err := host.Info()
	if err != nil {
		fmt.Print(err)
	}

	return Host{
		Name:     host.Hostname,
		OS:       host.OS,
		Arch:     runtime.GOARCH,
		Platform: host.Platform + " " + host.PlatformVersion,
		Uptime:   host.Uptime,
	}
}

func HostInfo(c *fiber.Ctx) {
	c.JSON(CheckHost())
}
