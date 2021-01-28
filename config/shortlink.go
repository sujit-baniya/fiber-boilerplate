package config

import (
	"fmt"
	"github.com/sujit-baniya/shorten"
)

type ShortlinkConfig struct {
	RedisUri   string       `yaml:"redis_uri" env-default:"redis://localhost:6379/0"`
	LinkLength int          `yaml:"link_length" env-default:"8"`
	EnableCSP  bool         `yaml:"enable_csp" env-default:"true"`
	LinkType   string       `yaml:"link_type" env-default:"l"`
	Server     ServerConfig `yaml:"server"`
	LinkOption shorten.LinkOption
}

func (v *ShortlinkConfig) Load() {
	shorten.LoadRedis(v.RedisUri)
	v.LinkOption = shorten.LinkOption{
		LinkLength: v.LinkLength,
		EnableCSP:  v.EnableCSP,
		Type:       v.LinkType,
		Host:       fmt.Sprintf("%s:%s", v.Server.Host, v.Server.Port),
		ExpiresIn:  0,
	}
}
