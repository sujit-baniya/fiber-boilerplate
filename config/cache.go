package config

import (
	"github.com/go-redis/redis"
	. "github.com/itsursujit/fiber-boilerplate/app"
	"log"
)

type CacheConfiguration struct {
	Cache_DSN string
	Cache_DB  int
}

var CacheConfig *CacheConfiguration //nolint:gochecknoglobals

func LoadCacheConfig() {
	loadDefaultCacheConfig()
	ViperConfig.Unmarshal(&CacheConfig)
	option, err := redis.ParseURL("redis://127.0.0.1:6379/0")
	if err != nil {
		log.Fatal(err)
	}
	RedisClient = redis.NewClient(option)
}

func loadDefaultCacheConfig() {
	ViperConfig.SetDefault("CACHE_DSN", "127.0.0.1:6379")
	ViperConfig.SetDefault("CACHE_DB", 0)
}
