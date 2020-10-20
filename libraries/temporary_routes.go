package libraries

import (
	fmt "fmt"
	"github.com/gofiber/fiber/v2" //nolint:goimports
	"time"
)

type URL struct {
	App *fiber.App
	Ctx *fiber.Ctx
}

type URLConfig struct {
	Expire    time.Duration //nolint:gofmt
	Lookup    func(interface{}) string
	AllowOnce bool
}

func (u *URL) TemporaryUrl(path string, config URLConfig) string {
	if !u.validateExistingRoute(path) {
		panic(fmt.Errorf("%s not found", path))
	}
	return ""
}

func (u *URL) SignedUrl(path string, config URLConfig) string {

	if !u.validateExistingRoute(path) {
		panic(fmt.Errorf("%s not found", path))
	}
	return ""
}

func (u *URL) TemporarySignedUrl(path string, config URLConfig) string {
	if !u.validateExistingRoute(path) {
		panic(fmt.Errorf("%s not found", path))
	}
	return ""
}

func (u *URL) validateExistingRoute(path string) bool {
	routes := u.App.Stack()
	for _, paths := range routes {
		for _, r := range paths {
			if r.Path == path {
				return true
			}
		}
	}
	return false
}
