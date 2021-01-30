package config

import (
	"github.com/gofiber/storage/redis"
)

type CacheConfig struct {
	*redis.Storage
	Driver string `yaml:"driver" env:"CACHE_DRIVER"`
	Name   string `yaml:"name" env:"CACHE_NAME"`
	Host   string `yaml:"host" env:"CACHE_HOST"`
	Port   int    `yaml:"port" env:"CACHE_PORT"`
	DB     int    `yaml:"db" env:"CACHE_DB"`
}

func (c *CacheConfig) Setup() {
	switch c.Driver {
	case "memcache":
	default:
		// Initialize custom config
		store := redis.New(redis.Config{
			Host:     c.Host,
			Port:     c.Port,
			Database: c.DB,
		})
		c.Storage = store
	}

}
