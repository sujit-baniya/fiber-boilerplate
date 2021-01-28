package config

import "github.com/gofiber/storage/redis"

type StorageConfig struct {
	*redis.Storage
	Driver   string `yaml:"driver" env:"STORAGE_DRIVER"`
	Name     string `yaml:"name" env:"STORAGE_NAME"`
	Host     string `yaml:"host" env:"STORAGE_HOST"`
	Username string `yaml:"username" env:"STORAGE_USERNAME"`
	Password string `yaml:"password" env:"STORAGE_PASSWORD"`
	Port     int    `yaml:"port" env:"STORAGE_PORT"`
	DB       int    `yaml:"db" env:"STORAGE_DB"`
}

func (c *StorageConfig) Setup() {
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
