package providers

import "github.com/thomasvvugt/fiber-boilerplate/app/configuration"

var appConfig *configuration.Configuration

func SetConfiguration(config *configuration.Configuration)  {
	appConfig = config
}

func GetConfiguration() (config *configuration.Configuration) {
	return appConfig
}
