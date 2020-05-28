package configuration

import (
	"github.com/gofiber/cors"
	"github.com/spf13/viper"
)

func loadCORSConfiguration() (enabled bool, config cors.Config, err error) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("cors")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultCORSConfiguration(provider)

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

	// Unmarshal the configuration file into logger.Config
	err = provider.Unmarshal(&config)

	// Return the configuration (and error if occurred)
	return provider.GetBool("Enabled"), config, err
}

// Set default configuration for the Logger Middleware
func setDefaultCORSConfiguration(provider *viper.Viper) {
	provider.SetDefault("Enabled", true)
	provider.SetDefault("AllowOrigins", []string{"*"})
	provider.SetDefault("AllowMethods", []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"})
	provider.SetDefault("AllowCredentials", nil)
	provider.SetDefault("ExposeHeaders", nil)
	provider.SetDefault("MaxAge", 0)
}
