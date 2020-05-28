package configuration

import (
	"log"
	"strings"
	"time"

	fsession "github.com/fasthttp/session/v2"
	"github.com/fasthttp/session/v2/providers/memory"
	"github.com/gofiber/session"
	"github.com/gofiber/session/provider/redis"

	"github.com/spf13/viper"
)

func loadSessionConfiguration() (enabled bool, config session.Config, err error) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("session")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultSessionConfiguration(provider)

	// Read configuration file
	err = provider.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error since we have default configurations
		} else {
			// Config file was found but another error was produced
			return provider.GetBool("Enabled"), config, err
		}
	}

	// Unmarshal the configuration file into session.Config
	err = provider.Unmarshal(&config)

	var driver fsession.Provider
	switch strings.ToLower(provider.GetString("Driver")) {
		case "redis":
			// Unmarshal the configuration file into redis.Config
			var redisConfig redis.Config
			err = provider.Unmarshal(&redisConfig)
			driver = redis.New(redisConfig)
	}

	// Set the provider depending on the Driver
	config.Provider = driver

	// Return the configuration (and error if occurred)
	return provider.GetBool("Enabled"), config, err
}

// Set default configuration for the Session Middleware
func setDefaultSessionConfiguration(provider *viper.Viper) {
	memoryProvider, err := memory.New(memory.Config{})
	if err != nil {
		log.Fatalf("Error when defining the default memory session provider: %v", err)
	}
	provider.SetDefault("Enabled", true)
	provider.SetDefault("Lookup", "cookie:session_id")
	provider.SetDefault("Secure", false)
	provider.SetDefault("Domain", "")
	provider.SetDefault("Expiration", 12 * time.Hour)
	provider.SetDefault("Provider", memoryProvider)
	provider.SetDefault("Generator", nil)
	provider.SetDefault("GCInterval", 1 * time.Minute)
}
