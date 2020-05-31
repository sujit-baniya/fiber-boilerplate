package web

import (
	"fmt"
	"github.com/gofiber/fiber"
	"strings"

	"github.com/thomasvvugt/fiber-boilerplate/app/providers"
)

func ShowLoginForm(c *fiber.Ctx) {
	store := providers.SessionProvider().Get(c)
	userid := store.Get("userid")
	if userid != nil {
		c.Redirect("/home")
		return
	}
	if err := c.Render("login", fiber.Map{}); err != nil {
		c.Status(500).Send(err.Error())
	}
}

func PostLoginForm(c *fiber.Ctx) {
	username := c.FormValue("username")
	// Find user
	user, err := FindUserByUsername(username)
	if err != nil {
		c.Redirect("/login")
		fmt.Printf("Error when finding user 1: %v", err)
		return
	}
	// Check if password matches hash
	if providers.HashProvider() != nil {
		password := c.FormValue("password")
		match, err := providers.HashProvider().MatchHash(password, user.Password)
		if err != nil {
			c.Redirect("/login")
			fmt.Printf("Error when finding user 2: %v", err)
		}
		if match {
			store := providers.SessionProvider().Get(c)
			defer store.Save()
			// Set the user ID in the session store
			store.Set("userid", user.ID)
			c.Send("You should be logged in successfully!")
		} else {
			c.Send("The entered details do not match our records.")
		}
	} else {
		c.Redirect("/login")
		fmt.Printf("Error Hash provider not set: %v", err)
		return
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
			c.Set("Set-Cookie", split[1]+"=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; HttpOnly")
			c.Send("You are now logged out.")
			return
		}
	}
}
