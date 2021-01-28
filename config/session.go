package config

import (
	"fmt"
	"time"

	"github.com/gofiber/session/v2"
	"github.com/gofiber/session/v2/provider/redis"
)

type SessionConfig struct {
	*session.Session
	Driver string `yaml:"driver" env:"SESSION_DRIVER"`
	Name   string `yaml:"name" env:"SESSION_NAME"`
	Host   string `yaml:"host" env:"SESSION_HOST"`
	Port   int    `yaml:"port" env:"SESSION_PORT"`
	DB     int    `yaml:"db" env:"SESSION_DB"`
}

func (s *SessionConfig) Setup() error {
	provider, err := redis.New(redis.Config{
		KeyPrefix:   "verify_rest_",
		Addr:        fmt.Sprintf("%s:%d", s.Host, s.Port),
		PoolSize:    8,                //nolint:gomnd
		IdleTimeout: 30 * time.Second, //nolint:gomnd
		DB:          s.DB,
		/*MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,*/
	})
	if err != nil {
		return err
	}
	s.Session = session.New(session.Config{
		Provider: provider,
	})
	return nil

}
