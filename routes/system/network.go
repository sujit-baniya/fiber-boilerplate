package system

import (
	"fmt"
	"net"
	"strings"

	"github.com/gofiber/fiber"
)

type address struct {
	IP string `json:"ip"`
}

type networkDevice struct {
	Name      string    `json:"name"`
	Addresses []address `json:"addresses"`
	MAC       string    `json:"mac"`
	Active    bool      `json:"active"`
}

type NetworkDevices []networkDevice

func CheckNetwork() []networkDevice {
	var networkDevices []networkDevice
	var mac string
	interfaces, err := net.Interfaces()

	if err != nil {
		fmt.Print(err)
	}

	for _, iface := range interfaces {
		byNameInterface, err := net.InterfaceByName(iface.Name)
		if err != nil {
			fmt.Println(err)
		}

		if iface.HardwareAddr.String() == "" {
			mac = "00:00:00:00:00:00"
		} else {
			mac = iface.HardwareAddr.String()
		}

		networkDevice := networkDevice{
			Name: iface.Name,
			MAC:  mac,
		}

		if (iface.Flags & net.FlagUp) == 0 {
			networkDevice.Active = false
		} else {
			networkDevice.Active = true
		}

		addresses, err := byNameInterface.Addrs()

		if err != nil {
			fmt.Println(err)
		}

		var addr []address

		for _, addrss := range addresses {
			addr = append(addr, address{
				IP: strings.Split(addrss.String(), "/")[0],
			})
		}

		networkDevice.Addresses = addr

		networkDevices = append(networkDevices, networkDevice)
	}

	return networkDevices
}

func NetworkInfo(c *fiber.Ctx) {
	c.JSON(CheckNetwork())
}
