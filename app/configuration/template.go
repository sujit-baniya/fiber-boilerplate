package configuration
/*
import (
	"github.com/gofiber/fiber"
	"github.com/spf13/viper"
	"strings"
)


import (
	"github.com/eknkc/amber"
	"github.com/gofiber/fiber"
	"github.com/gofiber/template/handlebars"
	"github.com/gofiber/template/mustache"
	"github.com/gofiber/template/pug"
	"strings"

	"github.com/spf13/viper"
)

func loadTemplateConfiguration() (enabled bool, engine *fiber.Templates) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("template")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultTemplateConfiguration(provider)

	// Read configuration file
	_ = provider.ReadInConfig()

	// Go over the provided configuration
	switch strings.ToLower(provider.GetString("Engine")) {
		case "mustache":
			engine = mustache.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "amber":
			engine = amber.New()
		case "handlebars":
			engine = handlebars.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "pug":
			engine = pug.New(provider.GetString("Folder"), provider.GetString("Extension"))
		default:
			engine = nil
	}

	// Return the configuration
	return provider.GetBool("Enabled"), engine
}

// Set default configuration for the Template Middleware
func setDefaultTemplateConfiguration(provider *viper.Viper) {
	provider.SetDefault("Enabled", false)
	provider.SetDefault("Engine", nil)
	provider.SetDefault("Folder", "")
	provider.SetDefault("Extension", ".html")
}
*/