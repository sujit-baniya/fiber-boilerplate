package config

type QueueConfig struct {
	Driver string `yaml:"driver" env:"CACHE_DRIVER"`
	Name   string `yaml:"name" env:"CACHE_NAME"`
	Host   string `yaml:"host" env:"CACHE_HOST"`
	Port   string `yaml:"port" env:"CACHE_PORT"`
	DB     string `yaml:"db" env:"CACHE_DB"`
}
