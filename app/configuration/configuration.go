package configuration

import (
	"github.com/gofiber/compression"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"
	"github.com/gofiber/session"

	hashing "github.com/thomasvvugt/fiber-hashing"
)

type Configuration struct {
	Fiber fiber.Settings
	App ApplicationConfiguration
	Enabled map[string]bool
	Logger logger.Config
	Recover recover.Config
	Compression compression.Config
	CORS cors.Config
	Helmet helmet.Config
	Hash hashing.Config
	Session session.Config
	PublicPrefix string
	PublicRoot string
	Public fiber.Static
	Database DatabaseConfiguration
}

func LoadConfigurations() (config Configuration, err error) {
	config.Enabled = make(map[string]bool)
	// Load the Fiber application configuration
	fiberSettings, err := loadFiberConfiguration()
	if err != nil {
		return config, err
	}
	config.Fiber = fiberSettings

	// Load the application configuration
	appConfig, err := loadApplicationConfiguration()
	if err != nil {
		return config, err
	}
	config.App = appConfig

	// Load the logger middleware configuration
	loggerEnabled, loggerConfig, err := loadLoggerConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["logger"] = loggerEnabled
	config.Logger = loggerConfig

	// Load the recover middleware configuration
	recoverEnabled, recoverConfig, err := loadRecoverConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["recover"] = recoverEnabled
	config.Recover = recoverConfig

	// Load the compression middleware configuration
	compressionEnabled, compressionConfig, err := loadCompressionConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["compression"] = compressionEnabled
	config.Compression = compressionConfig

	// Load the CORS middleware configuration
	corsEnabled, corsConfig, err := loadCORSConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["cors"] = corsEnabled
	config.CORS = corsConfig

	// Load the Helmet middleware configuration
	helmetEnabled, helmetConfig, err := loadHelmetConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["helmet"] = helmetEnabled
	config.Helmet = helmetConfig

	// Load the hashing configuration
	hashEnabled, hashConfig, err := loadHashConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["hash"] = hashEnabled
	config.Hash = hashConfig

	// Load the session middleware configuration
	sessionEnabled, sessionConfig, err := loadSessionConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["session"] = sessionEnabled
	config.Session = sessionConfig

	// Load the public, static files configuration
	publicEnabled, publicPrefix, publicRoot, publicConfig, err := loadPublicConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["public"] = publicEnabled
	config.PublicPrefix = publicPrefix
	config.PublicRoot = publicRoot
	config.Public = publicConfig

	// Load the database configuration
	databaseEnabled, databaseConfig, err := loadDatabaseConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["database"] = databaseEnabled
	config.Database = databaseConfig

	// Return the configuration
	return config, nil
}
