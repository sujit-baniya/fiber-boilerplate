package configuration

import (
	"github.com/spf13/viper"

	hashing "github.com/thomasvvugt/fiber-hashing"
	"github.com/thomasvvugt/fiber-hashing/driver/argon2id"
)

func loadHashConfiguration() (enabled bool, config hashing.Config, err error) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("hash")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultHashConfiguration(provider)

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

// Set default configuration for the Session Middleware
func setDefaultHashConfiguration(provider *viper.Viper) {
	provider.SetDefault("Enabled", true)
	provider.SetDefault("Driver", argon2id.New())
}
