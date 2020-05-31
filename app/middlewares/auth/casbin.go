package auth

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2" //nolint:goimports
	"github.com/gofiber/fiber"
	"github.com/thomasvvugt/fiber-boilerplate/app/configuration"
	"strings"
)

var authAdapter *gormadapter.Adapter
var enforcer *casbin.Enforcer

// Config holds the configuration for the middleware
type Config struct {
	// ModelFilePath is path to model file for Casbin.
	// Optional. Default: "./model.conf".
	ModelFilePath string

	// PolicyAdapter is an interface for different persistent providers.
	// Optional. Default: fileadapter.NewAdapter("./policy.csv").
	PolicyAdapter *gormadapter.Adapter

	// Lookup is a function that is used to look up current subject.
	// An empty string is considered as unauthenticated user.
	// Optional. Default: func(c *fiber.Ctx) string { return "" }
	Lookup func(*fiber.Ctx) string

	// Unauthorized defines the response body for unauthorized responses.
	// Optional. Default: func(c *fiber.Ctx) string { c.SendStatus(401) }
	Unauthorized func(*fiber.Ctx)

	// Forbidden defines the response body for forbidden responses.
	// Optional. Default: func(c *fiber.Ctx) string { c.SendStatus(403) }
	Forbidden func(*fiber.Ctx)
}

// CasbinMiddleware ...
type CasbinMiddleware struct {
	config   Config
	enforcer *casbin.Enforcer
}

// New creates an authorization middleware for use in Fiber
func New(config ...Config) *CasbinMiddleware {

	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.ModelFilePath == "" {
		cfg.ModelFilePath = "./model.conf"
	}

	if cfg.Lookup == nil {
		cfg.Lookup = func(c *fiber.Ctx) string { return "" }
	}

	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) {
			c.Redirect("/401")
			return
		}
	}

	if cfg.Forbidden == nil {
		cfg.Forbidden = func(c *fiber.Ctx) {
			c.Redirect("/403")
			return
		}
	}

	enforcer, err := casbin.NewEnforcer(cfg.ModelFilePath, cfg.PolicyAdapter)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database %s", err))
	}
	return &CasbinMiddleware{
		config:   cfg,
		enforcer: enforcer,
	}
}

type validationRule int

const (
	matchAll validationRule = iota
	atLeastOne
)

// MatchAll is an option that defines all permissions
// or roles should match the user.
var MatchAll = func(o *Options) {
	o.ValidationRule = matchAll
}

// AtLeastOne is an option that defines at least on of
// permissions or roles should match to pass.
var AtLeastOne = func(o *Options) {
	o.ValidationRule = atLeastOne
}

// PermissionParserFunc is used for parsing the permission
// to extract object and action usually
type PermissionParserFunc func(str string) []string

func permissionParserWithSeperator(sep string) PermissionParserFunc {
	return func(str string) []string {
		return strings.Split(str, sep)
	}
}

// PermissionParserWithSeperator is an option that parses permission
// with seperators
func PermissionParserWithSeperator(sep string) func(o *Options) {
	return func(o *Options) {
		o.PermissionParser = permissionParserWithSeperator(sep)
	}
}

// Options holds options of middleware
type Options struct {
	ValidationRule   validationRule
	PermissionParser PermissionParserFunc
}

// RequiresPermissions tries to find the current subject and determine if the
// subject has the required permissions according to predefined Casbin policies.
func (cm *CasbinMiddleware) RequiresPermissions(permissions []string, opts ...func(o *Options)) func(*fiber.Ctx) {

	options := &Options{
		ValidationRule:   matchAll,
		PermissionParser: permissionParserWithSeperator(":"),
	}

	for _, o := range opts {
		o(options)
	}

	return func(c *fiber.Ctx) {
		if len(permissions) == 0 {
			c.Next()
			return
		}

		sub := cm.config.Lookup(c)
		if len(sub) == 0 {
			cm.config.Unauthorized(c)
			return
		}

		if options.ValidationRule == matchAll {
			for _, permission := range permissions {
				vals := append([]string{sub}, options.PermissionParser(permission)...)
				if ok, err := cm.enforcer.Enforce(convertToInterface(vals)...); err != nil {
					c.SendStatus(fiber.StatusInternalServerError)
					return
				} else if !ok {
					cm.config.Forbidden(c)
					return
				}
			}
			c.Next()
			return
		} else if options.ValidationRule == atLeastOne {
			for _, permission := range permissions {
				vals := append([]string{sub}, options.PermissionParser(permission)...)
				if ok, err := cm.enforcer.Enforce(convertToInterface(vals)...); err != nil {
					c.SendStatus(fiber.StatusInternalServerError)
					return
				} else if ok {
					c.Next()
					return
				}
			}
			cm.config.Forbidden(c)
			return
		}

		c.Next()
	}
}

// RoutePermission tries to find the current subject and determine if the
// subject has the required permissions according to predefined Casbin policies.
// This method uses http Path and Method as object and action.
func (cm *CasbinMiddleware) RoutePermission() func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {

		sub := cm.config.Lookup(c)
		if len(sub) == 0 {
			cm.config.Unauthorized(c)
			return
		}

		if ok, err := cm.enforcer.Enforce(sub, c.Path(), c.Method()); err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return
		} else if !ok {
			cm.config.Forbidden(c)
			return
		}

		c.Next()
		return
	}
}

// RequiresRoles tries to find the current subject and determine if the
// subject has the required roles according to predefined Casbin policies.
func (cm *CasbinMiddleware) RequiresRoles(roles []string, opts ...func(o *Options)) func(*fiber.Ctx) {
	options := &Options{
		ValidationRule:   matchAll,
		PermissionParser: permissionParserWithSeperator(":"),
	}

	for _, o := range opts {
		o(options)
	}

	return func(c *fiber.Ctx) {
		if len(roles) == 0 {
			c.Next()
			return
		}

		sub := cm.config.Lookup(c)
		if len(sub) == 0 {
			cm.config.Unauthorized(c)
			return
		}

		userRoles, err := cm.enforcer.GetRolesForUser(sub)
		if err != nil {
			c.SendStatus(fiber.StatusInternalServerError)
			return
		}

		if options.ValidationRule == matchAll {
			for _, role := range roles {
				if !contains(userRoles, role) {
					cm.config.Forbidden(c)
					return
				}
			}
			c.Next()
			return
		} else if options.ValidationRule == atLeastOne {
			for _, role := range roles {
				if contains(userRoles, role) {
					c.Next()
					return
				}
			}
			cm.config.Forbidden(c)
			return
		}

		c.Next()
	}
}

func contains(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func convertToInterface(arr []string) []interface{} {
	in := make([]interface{}, 0)
	for _, a := range arr {
		in = append(in, a)
	}
	return in
}

func CasbinRbac(config *configuration.DatabaseConfiguration) *gormadapter.Adapter {
	var err error
	if authAdapter != nil {
		return authAdapter
	}

	switch strings.ToLower(config.Driver) { //nolint:wsl
	case "mssql":
		authAdapter, err = gormadapter.NewAdapter("mssql", "sqlserver://"+config.Username+":"+config.Password+"@"+config.Host+":"+string(config.Port)+"?database="+config.Database, true)
	case "mysql", "mariadb":
		authAdapter, err = gormadapter.NewAdapter("mysql", config.Username+":"+config.Password+"@tcp("+config.Host+")/"+config.Database+"?charset=utf8&parseTime=True&loc=Local", true)
	case "postgre", "postgres", "postgresql":
		authAdapter, err = gormadapter.NewAdapter("postgres", "host="+config.Host+" port="+string(config.Port)+" user="+config.Username+" dbname="+config.Database+" password="+config.Password, true)
	case "sqlite", "sqlite3":
		authAdapter, err = gormadapter.NewAdapter("sqlite3", config.Database)
	}
	// enforcer, err := casbin.NewEnforcer("./config/rbac_model.conf", authAdapter)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database %s", err))
	}
	return authAdapter
}
