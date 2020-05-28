package configuration

import (
	"github.com/gofiber/helmet"
	"github.com/spf13/viper"
)

func loadHelmetConfiguration() (enabled bool, config helmet.Config, err error) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("helmet")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultHelmetConfiguration(provider)

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
func setDefaultHelmetConfiguration(provider *viper.Viper) {
	provider.SetDefault("Enabled", true)
	provider.SetDefault("XSSProtection", "1; mode=block")
	provider.SetDefault("ContentTypeNosniff", "nosniff")
	provider.SetDefault("XFrameOptions", "SAMEORIGIN")
	provider.SetDefault("HSTSMaxAge", 31536000)
	provider.SetDefault("HSTSExcludeSubdomains", false)
	provider.SetDefault("ContentSecurityPolicy", "")
	provider.SetDefault("CSPReportOnly", false)
	provider.SetDefault("HSTSPreloadEnabled", false)
	provider.SetDefault("ReferrerPolicy", "")
}
