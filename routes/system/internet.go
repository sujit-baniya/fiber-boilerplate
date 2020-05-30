package system

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gofiber/fiber"
)

type Setting struct {
	IP  string `xml:"ip,attr" json:"ip"`
	Lat string `xml:"lat,attr" json:"lat"`
	Lon string `xml:"lon,attr" json:"lon"`
	Isp string `xml:"isp,attr" json:"isp"`
}

type Settings struct {
	Settings []Setting `xml:"client"`
}

func checkError(e error) {
	if e != nil {
		log.Println(e)
	}
}

func CheckInternet() Setting {
	// Fetch xml user data
	resp, err := http.Get("http://speedtest.net/speedtest-config.php")
	checkError(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	defer resp.Body.Close()

	// Decode xml
	decoder := xml.NewDecoder(bytes.NewReader(body))
	settings := Settings{}
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			decoder.DecodeElement(&settings, &se)
		}
	}
	if settings.Settings == nil {
		fmt.Println("Warning: Cannot fetch user information. http://www.speedtest.net/speedtest-config.php is temporarily unavailable.")
		return Setting{}
	}
	return settings.Settings[0]
}

func InternetInfo(c *fiber.Ctx) {
	c.JSON(CheckInternet())
}
