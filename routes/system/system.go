package system

import "github.com/gofiber/fiber"

type System struct {
	Host      Host                     `json:"host"`
	CPU       CPU                      `json:"cpu"`
	RAM       RAM                      `json:"ram"`
	Disks     Disks                    `json:"disks"`
	Bandwidth []NetworkDeviceBandwidth `json:"bandwidth"`
}

func CheckSystem() System {
	sys := System{
		Host:      CheckHost(),
		CPU:       CheckCpu(),
		RAM:       CheckRam(),
		Disks:     CheckDisk(),
		Bandwidth: CheckBandwidth(),
	}
	return sys
}

func SystemInfo(c *fiber.Ctx) {
	c.JSON(CheckSystem())
}
