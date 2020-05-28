package providers

import "github.com/gofiber/fiber"

type Config struct {
	Username string
}

type Provider struct {
	Config Config
}

var authP Provider

func AuthProvider() *Provider {
	return &authP
}

func SetAuthProvider(config ...Config) {
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}
	if cfg.Username == "" {
		cfg.Username = "username"
	}
	authP = Provider{Config: cfg}
}

func IsAuthenticated(c *fiber.Ctx) (authenticated bool) {
	store := SessionProvider().Get(c)
	// Get User ID from session store
	userID, correct := store.Get("userid").(int64)
	if !correct {
		userID = 0
	}
	auth := false
	if userID > 0 {
		auth = true
	}
	return auth
}
