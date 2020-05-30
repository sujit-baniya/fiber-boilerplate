package system

import (
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/prometheus/procfs"
	"log"
	"time"
)

type NetworkDeviceBandwidth struct {
	Name      string `json:"name"`
	RxBytes   uint64 `json:"rxBytes"`
	TxBytes   uint64 `json:"txBytes"`
	RxPackets uint64 `json:"rxPackets"`
	TxPackets uint64 `json:"txPackets"`
}

func CheckBandwidth() []NetworkDeviceBandwidth {
	p, err := procfs.NewDefaultFS()
	if err != nil {
		log.Fatalf("could not get process: %s", err)
	}
	net, err := p.NetDev()
	if err != nil {
		fmt.Println(err)
	}

	networks := CheckNetwork()

	//Round 1
	var stats1 []NetworkDeviceBandwidth
	for _, netw := range networks {
		stats1 = append(stats1, NetworkDeviceBandwidth{
			Name:      netw.Name,
			RxBytes:   net[netw.Name].RxBytes,
			TxBytes:   net[netw.Name].TxBytes,
			RxPackets: net[netw.Name].RxPackets,
			TxPackets: net[netw.Name].TxPackets,
		})
	}

	time.Sleep(1000 * time.Millisecond)

	net, err = p.NetDev()
	if err != nil {
		fmt.Println(err)
	}

	//Round 2
	var stats2 []NetworkDeviceBandwidth
	for _, netw := range networks {
		stats2 = append(stats2, NetworkDeviceBandwidth{
			Name:      netw.Name,
			RxBytes:   net[netw.Name].RxBytes,
			TxBytes:   net[netw.Name].TxBytes,
			RxPackets: net[netw.Name].RxPackets,
			TxPackets: net[netw.Name].TxPackets,
		})
	}

	//DIFF
	var diffStats []NetworkDeviceBandwidth

	for i, netw := range networks {
		diffStats = append(diffStats, NetworkDeviceBandwidth{
			Name:      netw.Name,
			RxBytes:   stats2[i].RxBytes - stats1[i].RxBytes,
			TxBytes:   stats2[i].TxBytes - stats1[i].TxBytes,
			RxPackets: stats2[i].RxPackets,
			TxPackets: stats2[i].TxPackets,
		})
	}

	return diffStats
}
func BandwidthInfo(c *fiber.Ctx) {
	c.JSON(CheckBandwidth())

}
