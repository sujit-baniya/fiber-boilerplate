package configuration

import (
	"github.com/spf13/viper"
)

type ApplicationConfiguration struct {
	Listen interface{}
	SuppressWWW bool
	ForceHTTPS bool
}

func loadApplicationConfiguration() (ApplicationConfiguration, error) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("app")
	provider.AddConfigPath("./config")

	// Create a new fiber.Settings variable
	var config ApplicationConfiguration

	// Set default configurations
	setDefaultApplicationConfiguration(provider)

	// Read configuration file
	err := provider.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error since we have default configurations
		} else {
			// Config file was found but another error was produced
			return config, err
		}
	}

	// Unmarshal the configuration file into the ApplicationConfiguration struct
	err = provider.Unmarshal(&config)

	// Return the configuration (and error if occurred)
	return config, err
}

// Set default configuration for the application
func setDefaultApplicationConfiguration(provider *viper.Viper) {
	provider.SetDefault("Listen", "8080")
	provider.SetDefault("SuppressWWW", true)
	provider.SetDefault("ForceHTTPS", false)
}
