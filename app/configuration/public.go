package configuration

import (
	"github.com/gofiber/fiber"
	"github.com/spf13/viper"
)

func loadPublicConfiguration() (enabled bool, prefix string, root string, config fiber.Static, err error) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("public")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultPublicConfiguration(provider)

	// Read configuration file
	err = provider.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error since we have default configurations
		} else {
			// Config file was found but another error was produced
			return provider.GetBool("Enabled"), provider.GetString("Prefix"), provider.GetString("Root"), config, err
		}
	}

	// Unmarshal the configuration file into logger.Config
	err = provider.Unmarshal(&config)

	// Return the configuration (and error if occurred)
	return provider.GetBool("Enabled"), provider.GetString("Prefix"), provider.GetString("Root"), config, err
}

// Set default configuration for the Logger Middleware
func setDefaultPublicConfiguration(provider *viper.Viper) {
	provider.SetDefault("Enabled", true)
	provider.SetDefault("Prefix", "/")
	provider.SetDefault("Root", "./public")
	provider.SetDefault("Compress", false)
	provider.SetDefault("ByteRange", false)
	provider.SetDefault("Browse", false)
	provider.SetDefault("Index", "index.html")
}
