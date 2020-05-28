package web

import (
	"fmt"
	"github.com/gofiber/fiber"
	"log"
	"strings"

	"github.com/thomasvvugt/fiber-boilerplate/app/providers"
)

func ShowLoginForm(c *fiber.Ctx) {
	if err := c.Render("login", fiber.Map{}); err != nil {
		c.Status(500).Send(err.Error())
	}
}

func PostLoginForm(c *fiber.Ctx) {
	username := c.FormValue("username")
	// Find user
	user, err := FindUserByUsername(username)
	if err != nil {
		log.Fatalf("Error when finding user: %v", err)
	}
	// Check if password matches hash
	if providers.HashProvider() != nil {
		password := c.FormValue("password")
		match, err := providers.HashProvider().MatchHash(password, user.Password)
		if err != nil {
			log.Fatalf("Error when matching hash for password: %v", err)
		}
		if match {
			store := providers.SessionProvider().Get(c)
			defer store.Save()
			// Set the user ID in the session store
			store.Set("userid", user.ID)
			fmt.Printf("User set in session store with ID: %v\n", user.ID)
			c.Send("You should be logged in successfully!")
		} else {
			c.Send("The entered details do not match our records.")
		}
	} else {
		panic("Hash provider was not set")
	}
}

func PostLogoutForm(c *fiber.Ctx) {
	if providers.IsAuthenticated(c) {
		store := providers.SessionProvider().Get(c)
		store.Delete("userid")
		store.Save()
		// Check if cookie needs to be unset
		config := providers.GetConfiguration()
		lookup := config.Session.Lookup
		split := strings.Split(lookup, ":")
		if strings.ToLower(split[0]) == "cookie" {
			// Unset cookie on client-side
			c.Set("Set-Cookie", split[1] + "=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; HttpOnly")
			c.Send("You are now logged out.")
			return
		}
	}
}
