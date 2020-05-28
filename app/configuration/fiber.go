package configuration

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/template/amber"
	"github.com/gofiber/template/handlebars"
	"github.com/gofiber/template/html"
	"github.com/gofiber/template/mustache"
	"github.com/gofiber/template/pug"
	"github.com/spf13/viper"
	"strings"
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
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { //nolint:gosimple
			// Config file not found; ignore error since we have default configurations
		} else {
			// Config file was found but another error was produced
			return settings, err
		}
	}

	// Unmarshal the configuration file into fiber.Settings
	err = provider.Unmarshal(&settings)
	setTemplate(&settings)
	// Return the configuration (and error if occurred)
	return settings, err
}

func setTemplate(settings *fiber.Settings) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("template")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultTemplateConfiguration(provider)

	// Read configuration file
	_ = provider.ReadInConfig()

	if provider.GetBool("Enabled") {
		// Go over the provided configuration
		switch strings.ToLower(provider.GetString("Engine")) {
		case "mustache":
			settings.Templates = mustache.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "handlebars":
			settings.Templates = handlebars.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "pug":
			settings.Templates = pug.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "html":
			settings.Templates = html.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "amber":
			settings.Templates = amber.New(provider.GetString("Folder"), provider.GetString("Extension"))
		default:
			settings.Templates = html.New(provider.GetString("Folder"), provider.GetString("Extension"))
		}
	}
}

// Set default configuration for the Template Middleware
func setDefaultTemplateConfiguration(provider *viper.Viper) {
	provider.SetDefault("Enabled", false)
	provider.SetDefault("Engine", "html")
	provider.SetDefault("Folder", "./resources/views")
	provider.SetDefault("Extension", ".html")
}

// Set default configuration for Fiber
func setDefaultFiberConfiguration(provider *viper.Viper) {
	provider.SetDefault("Prefork", false)
	provider.SetDefault("ServerHeader", "")
	provider.SetDefault("StrictRouting", false)
	provider.SetDefault("CaseSensitive", false)
	provider.SetDefault("Immutable", false)
	provider.SetDefault("BodyLimit", 4 * 1024 * 1024)
	provider.SetDefault("Concurrency", 256 * 1024)
	provider.SetDefault("DisableKeepalive", false)
	provider.SetDefault("DisableDefaultDate", false)
	provider.SetDefault("DisableDefaultContentType", false)
	provider.SetDefault("DisableStartupMessage", false)
	provider.SetDefault("DisableKeepalive", false)
	provider.SetDefault("DisableDefaultDate", false)
	provider.SetDefault("DisableDefaultContentType", false)
	provider.SetDefault("DisableStartupMessage", false)
	provider.SetDefault("ETag", false)
	provider.SetDefault("ReadTimeout", nil)
	provider.SetDefault("WriteTimeout", nil)
	provider.SetDefault("IdleTimeout", nil)
}
