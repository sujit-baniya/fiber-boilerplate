package config

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type AuthConfig struct {
	*gormadapter.Adapter
	Type     string         `yaml:"type" env:"AUTH_TYPE" env-default:"simple"`
	DB       DatabaseDriver `yaml:"driver"`
	Casbin   *Casbin
	Enforcer *casbin.Enforcer
}

func (d *AuthConfig) Setup() {
	var err error
	connectionString := ""
	switch d.DB.Driver {
	case "postgres":
		connectionString = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", d.DB.Host, d.DB.Port, d.DB.Username, d.DB.DBName, d.DB.Password)
	default:
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", d.DB.Username, d.DB.Password, d.DB.Host, d.DB.Port, d.DB.DBName)
	}
	adapter, err := gormadapter.NewAdapter(d.DB.Driver, connectionString)

	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}
	d.Adapter = adapter
	enforcer, err := casbin.NewEnforcer("rbac_model.conf")
	if err != nil {
		panic(err)
	}
	enforcer.SetAdapter(adapter)
	_ = enforcer.LoadPolicy()
	d.Enforcer = enforcer
}
