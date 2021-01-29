package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sujit-baniya/flash"
	"github.com/sujit-baniya/ip"
	"github.com/sujit-baniya/log"
	"github.com/sujit-baniya/sblogger"
	"os"
	"path/filepath"
	"time"
)

// Config is a application configuration structure
type AppConfig struct {
	Auth       AuthConfig `yaml:"auth"`
	Mail       Mail       `yaml:"mail"`
	Hash       Hash
	View       ViewConfig      `yaml:"view"`
	Cache      CacheConfig     `yaml:"cache"`
	Database   DatabaseConfig  `yaml:"database"`
	Session    SessionConfig   `yaml:"session"`
	Queue      QueueConfig     `yaml:"queue"`
	PayPal     PayPalConfig    `yaml:"paypal"`
	JwtSecrets JwtSecrets      `yaml:"jwt"`
	Storage    StorageConfig   `yaml:"storage"`
	Server     ServerConfig    `yaml:"server"`
	Log        LogConfig       `yaml:"log"`
	Token      Token           `yaml:"token"`
	Shortlink  ShortlinkConfig `yaml:"shortlink"`
	Spam       SpamConfig      `yaml:"spam"`
	Profiler   ProfilerConfig  `yaml:"profiler"`
	Flash      *flash.Flash
	GeoIP      *ip.GeoIpDB
	ConfigFile string
}

func (cfg *AppConfig) Setup() {
	err := godotenv.Load()
	// read configuration from the file and environment variables
	if err = cleanenv.ReadConfig(cfg.ConfigFile, cfg); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	cfg.Server.LoadPath()
	cfg.View.Load(cfg.Server.Path)
	cfg.Mail.View = &cfg.View
	cfg.Server.TemplateEngine = cfg.View.Template.TemplateEngine
	cfg.Server.Setup()
	cfg.LoadComponents()
	if cfg.Auth.Type == "casbin" {
		cfg.Auth.Setup(cfg.Database.DB)
	}
	path := MakeDir(filepath.Join(cfg.Server.AssetPath, "GeoLite2-City.mmdb"))
	cfg.GeoIP = ip.NewGeoIpDB(path)
}

func (cfg *AppConfig) PrepareLog() {
	writer := &log.MultiWriter{}
	path := MakeDir(filepath.Join(cfg.Server.Path, cfg.Log.InfoLevel.Path))
	writer.InfoWriter = &log.FileWriter{Filename: filepath.Join(path, "INFO.log"), EnsureFolder: true, TimeFormat: cfg.Log.InfoLevel.TimeFormat}

	path = MakeDir(filepath.Join(cfg.Server.Path, cfg.Log.WarnLevel.Path))
	writer.WarnWriter = &log.FileWriter{Filename: filepath.Join(cfg.Server.Path, cfg.Log.WarnLevel.Path, "WARN.log"), EnsureFolder: true, TimeFormat: cfg.Log.WarnLevel.TimeFormat}

	path = MakeDir(filepath.Join(cfg.Server.Path, cfg.Log.ErrorLevel.Path))
	writer.ErrorWriter = &log.FileWriter{Filename: filepath.Join(cfg.Server.Path, cfg.Log.ErrorLevel.Path, "ERROR.log"), EnsureFolder: true, TimeFormat: cfg.Log.ErrorLevel.TimeFormat}
	if cfg.Log.ConsoleLog.Show {
		writer.ConsoleWriter = &log.IOWriter{Writer: os.Stderr}
		writer.ConsoleLevel = log.InfoLevel
	}
	log.DefaultLogger = log.Logger{
		TimeField:  cfg.Log.TimeField,
		TimeFormat: cfg.Log.TimeFormat,
		Writer:     writer,
	}
}

func (cfg *AppConfig) Route404() {
	cfg.Server.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Page not found")
	})
}

func (cfg *AppConfig) LoadComponents() {
	cfg.Flash = &flash.Flash{
		CookiePrefix: "Verify-Rest",
	}
	cfg.LoadStatic()
	cfg.PrepareLog()
	cfg.Server.Use(
		sblogger.New(sblogger.Config{
			Logger:    &log.DefaultLogger,
			LogWriter: log.DefaultLogger.Writer,
		}),
	)
	_ = cfg.Database.Setup()
	_ = cfg.Session.Setup()
	cfg.Cache.Setup()
	cfg.Storage.Setup()
	cfg.Shortlink.Load()
	cfg.Spam.Load(cfg.Server.AssetPath)
}

func (cfg *AppConfig) LoadStatic() {
	cfg.Server.Static("/websocket", "./resources/views/websocket.html")
	cfg.Server.Static("/", filepath.Join(cfg.Server.Path, cfg.Server.PublicPath), fiber.Static{
		Compress:      true,
		ByteRange:     true,
		CacheDuration: 24 * time.Hour,
	})
}

func (cfg *AppConfig) LoadSpamDetectionEngine() {

}
