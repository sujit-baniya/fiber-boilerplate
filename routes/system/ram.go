package system

import (
	"fmt"

	"github.com/gofiber/fiber"
	"github.com/shirou/gopsutil/mem"
)

type RAM struct {
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
	Usage uint64 `json:"usage"`
}

func CheckRam() RAM {
	memory, err := mem.VirtualMemory()
	if err != nil {
		fmt.Print(err)
	}
	ram := RAM{
		Free:  memory.Total - memory.Used,
		Usage: memory.Used,
		Total: memory.Total,
	}

	return ram
}

func RamInfo(c *fiber.Ctx) {
	c.JSON(CheckRam())
}
