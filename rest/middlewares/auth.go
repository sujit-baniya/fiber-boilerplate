package middlewares

import (
	"errors"
	"fmt"
	"github.com/sujit-baniya/fiber-boilerplate/app"
	config2 "github.com/sujit-baniya/fiber-boilerplate/config"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/auth"
	"github.com/sujit-baniya/log"
	"reflect"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// Config defines the config for BasicAuth middleware
type AuthConfig struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool

	// SuccessHandler defines a function which is executed for a valid token.
	// Optional. Default: nil
	SuccessHandler fiber.Handler

	// ErrorHandler defines a function which is executed for an invalid token.
	// It may be used to define a custom JWT error.
	// Optional. Default: 401 Invalid or expired JWT
	ErrorHandler fiber.ErrorHandler

	// Signing key to validate token. Used as fallback if SigningKeys has length 0.
	// Required. This or SigningKeys.
	SigningKey interface{}

	// Map of signing keys to validate token with kid field usage.
	// Required. This or SigningKey.
	SigningKeys map[string]interface{}

	// Signing method, used to check token signing method.
	// Optional. Default: "HS256".
	// Possible values: "HS256", "HS384", "HS512", "ES256", "ES384", "ES512", "RS256", "RS384", "RS512"
	SigningMethod string

	// Context key to store user information from the token into context.
	// Optional. Default: "user".
	ContextKey string

	// Claims are extendable claims data defining token content.
	// Optional. Default value jwt.MapClaims
	Claims jwt.Claims

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "param:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// AuthScheme to be used in the Authorization header.
	// Optional. Default: "Bearer".
	AuthScheme string

	keyFunc jwt.Keyfunc
}

// New ...
func Authenticate(config ...AuthConfig) fiber.Handler {
	// Init config
	var cfg AuthConfig
	if len(config) > 0 {
		cfg = config[0]
	}
	if cfg.SuccessHandler == nil {
		cfg.SuccessHandler = func(c *fiber.Ctx) error {
			return c.Next()
		}
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			var er fiber.Error
			if err.Error() == "Missing or malformed JWT" {
				er.Code = fiber.StatusBadRequest
			} else {
				er.Code = fiber.StatusUnauthorized
			}
			er.Message = err.Error()
			return config2.CustomErrorHandler(c, &er)
		}
	}
	if cfg.SigningKey == nil && len(cfg.SigningKeys) == 0 {
		log.Error().Msg("Fiber: JWT middleware requires signing key")
	}
	if cfg.SigningMethod == "" {
		cfg.SigningMethod = "HS256"
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = "user"
	}
	if cfg.Claims == nil {
		cfg.Claims = jwt.MapClaims{}
	}
	if cfg.TokenLookup == "" {
		cfg.TokenLookup = "header:" + fiber.HeaderAuthorization
	}
	if cfg.AuthScheme == "" {
		cfg.AuthScheme = "Bearer"
	}
	cfg.keyFunc = func(t *jwt.Token) (interface{}, error) {
		// Check the signing method
		if t.Method.Alg() != cfg.SigningMethod {
			return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
		}
		if len(cfg.SigningKeys) > 0 {
			if kid, ok := t.Header["kid"].(string); ok {
				if key, ok := cfg.SigningKeys[kid]; ok {
					return key, nil
				}
			}
			return nil, fmt.Errorf("Unexpected jwt key id=%v", t.Header["kid"])
		}
		return cfg.SigningKey, nil
	}
	// Initialize
	extractors := make([]func(c *fiber.Ctx) (string, error), 0)
	rootParts := strings.Split(cfg.TokenLookup, ",")
	for _, rootPart := range rootParts {
		parts := strings.Split(strings.TrimSpace(rootPart), ":")

		switch parts[0] {
		case "header":
			extractors = append(extractors, jwtFromHeader(parts[1], cfg.AuthScheme))
		case "query":
			extractors = append(extractors, jwtFromQuery(parts[1]))
		case "param":
			extractors = append(extractors, jwtFromParam(parts[1]))
		case "cookie":
			extractors = append(extractors, jwtFromCookie(parts[1]))
		}
	}
	// Return middleware handler
	return func(c *fiber.Ctx) error {
		// Filter request to skip middleware
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}
		var auth string
		var err error

		for _, extractor := range extractors {
			auth, err = extractor(c)
			if auth != "" && err == nil {
				break
			}
		}

		if err != nil {
			return cfg.ErrorHandler(c, err)
		}
		token := new(jwt.Token)
		if _, ok := cfg.Claims.(jwt.MapClaims); ok {
			token, err = jwt.Parse(auth, cfg.keyFunc)
		} else {
			t := reflect.ValueOf(cfg.Claims).Type().Elem()
			claims := reflect.New(t).Interface().(jwt.Claims)
			token, err = jwt.ParseWithClaims(auth, claims, cfg.keyFunc)
		}
		if err == nil && token.Valid {
			// Store user information from token into context.
			c.Locals(cfg.ContextKey, token)
			return cfg.SuccessHandler(c)
		}
		return cfg.ErrorHandler(c, err)
	}
}

// jwtFromHeader returns a function that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		auth := c.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", errors.New("Missing or malformed JWT")
	}
}

// jwtFromQuery returns a function that extracts token from the query string.
func jwtFromQuery(param string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Query(param)
		if token == "" {
			return "", errors.New("Missing or malformed JWT")
		}
		return token, nil
	}
}

// jwtFromParam returns a function that extracts token from the url param string.
func jwtFromParam(param string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Params(param)
		if token == "" {
			return "", errors.New("Missing or malformed JWT")
		}
		return token, nil
	}
}

// jwtFromCookie returns a function that extracts token from the named cookie.
func jwtFromCookie(name string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Cookies(name)
		if token == "" {
			return "", errors.New("Missing or malformed JWT")
		}
		return token, nil
	}
}

func AuthWeb() func(*fiber.Ctx) error {
	return Authenticate(AuthConfig{
		SigningKey:  []byte(app.Http.Token.AppJwtSecret),
		TokenLookup: "cookie:Verify-Rest-Token",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			auth.Logout(ctx)
			return ctx.Redirect("/login")
		},
	})
}

func AuthAdmin(c *fiber.Ctx) error {
	if !auth.IsAdmin(c) {
		auth.Logout(c)
		return c.Redirect("/login")
	}
	return c.Next()
}

func AuthApi() func(*fiber.Ctx) error {
	return Authenticate(AuthConfig{
		SigningKey:  []byte(app.Http.Token.ApiJwtSecret),
		TokenLookup: "header:Verify-Rest-Token",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			auth.Logout(ctx)
			return ctx.Status(401).JSON("Invalid Attempt")
		},
	})
}
