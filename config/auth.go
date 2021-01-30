package config

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthConfig struct {
	*gormadapter.Adapter
	Type     string `yaml:"type" env:"AUTH_TYPE" env-default:"simple"`
	Casbin   *Casbin
	Enforcer *casbin.Enforcer
}

func (d *AuthConfig) Setup(db *gorm.DB, file string) {
	adapter, err := gormadapter.NewAdapterByDB(db)

	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}
	d.Adapter = adapter
	enforcer, err := casbin.NewEnforcer(file)
	if err != nil {
		panic(err)
	}
	enforcer.SetAdapter(adapter)
	err = enforcer.LoadPolicy()
	if err != nil {
		panic(err)
	}
	d.Enforcer = enforcer
	authConf := CasbinAuthConfig{
		Enforcer:      d.Enforcer,
		PolicyAdapter: d.Adapter,
		Lookup: func(ctx *fiber.Ctx) string {
			return ctx.Locals("user_id").(string)
		},
		Unauthorized: func(c *fiber.Ctx) error {
			var err fiber.Error
			err.Code = fiber.StatusUnauthorized
			return CustomErrorHandler(c, &err)
		},
		Forbidden: func(c *fiber.Ctx) error {
			var err fiber.Error
			err.Code = fiber.StatusForbidden
			return CustomErrorHandler(c, &err)
		},
	}
	d.Casbin = CasbinAuth(authConf)
}
