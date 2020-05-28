package configuration

import (
	"github.com/gofiber/fiber"
	"os"

	"github.com/gofiber/recover"

	"github.com/spf13/viper"
)

func loadRecoverConfiguration() (enabled bool, config recover.Config, err error) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("recover")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultRecoverConfiguration(provider)

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

	// Unmarshal the configuration file into recover.Config
	err = provider.Unmarshal(&config)

	// Return the configuration (and error if occurred)
	return provider.GetBool("Enabled"), config, err
}

// Set default configuration for the Recover Middleware
func setDefaultRecoverConfiguration(provider *viper.Viper) {
	provider.SetDefault("Enabled", true)
	provider.SetDefault("Filter", nil)
	provider.SetDefault("Handler", func(c *fiber.Ctx, err error) { c.SendStatus(500); if err := c.Render("errors/500", fiber.Map{}); err != nil { c.Status(500).Send(err.Error()) } })
	provider.SetDefault("Log", false)
	provider.SetDefault("Output", os.Stderr)
}
