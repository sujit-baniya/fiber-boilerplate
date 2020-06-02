package configuration

import (
	"github.com/gofiber/fiber"

	"github.com/spf13/viper"
)

func loadFiberConfiguration() (settings fiber.Settings, err error) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("fiber")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultFiberConfiguration(provider)

	// Read configuration file
	err = provider.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error since we have default configurations
		} else {
			// Config file was found but another error was produced
			return settings, err
		}
	}

	// Unmarshal the configuration file into fiber.Settings
	err = provider.Unmarshal(&settings)

	// Return the configuration (and error if occurred)
	return settings, err
}

// Set default configuration for Fiber
func setDefaultFiberConfiguration(provider *viper.Viper) {
	provider.SetDefault("Prefork", false)
	provider.SetDefault("ServerHeader", "")
	provider.SetDefault("StrictRouting", false)
	provider.SetDefault("CaseSensitive", false)
	provider.SetDefault("Immutable", false)
	provider.SetDefault("BodyLimit", 4*1024*1024)
	provider.SetDefault("Concurrency", 256*1024)
	provider.SetDefault("DisableKeepalive", false)
	provider.SetDefault("DisableDefaultDate", false)
	provider.SetDefault("DisableDefaultContentType", false)
	provider.SetDefault("DisableStartupMessage", false)
	provider.SetDefault("ETag", false)
	provider.SetDefault("ReadTimeout", nil)
	provider.SetDefault("WriteTimeout", nil)
	provider.SetDefault("IdleTimeout", nil)
}
