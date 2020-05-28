package configuration

import (
	"github.com/spf13/viper"
)

type DatabaseConfiguration struct {
	Driver string
	Host string
	Port int
	Username string
	Password string
	Database string
}

func loadDatabaseConfiguration() (enabled bool, config DatabaseConfiguration, err error) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("database")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultDatabaseConfiguration(provider)

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

	// Unmarshal the configuration file into DatabaseConfiguration
	err = provider.Unmarshal(&config)

	// Return the configuration (and error if occurred)
	return provider.GetBool("Enabled"), config, err
}

// Set default configuration for the Database
func setDefaultDatabaseConfiguration(provider *viper.Viper) {
	provider.SetDefault("Enabled", true)
	provider.SetDefault("Driver", "mysql")
	provider.SetDefault("Host", "127.0.0.1")
	provider.SetDefault("Port", "3306")
	provider.SetDefault("Username", "fiber")
	provider.SetDefault("Password", "secret")
	provider.SetDefault("Database", "fiber")
}
