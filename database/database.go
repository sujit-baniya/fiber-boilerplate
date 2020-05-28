package database

import (
	"github.com/thomasvvugt/fiber-boilerplate/app/configuration"

	"strings"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func Instance() *gorm.DB {
	return db
}

func Connect(config *configuration.DatabaseConfiguration) {
	var err error
	switch strings.ToLower(config.Driver) {
	case "mssql":
		db, err = gorm.Open("mssql", "sqlserver://" + config.Username + ":" + config.Password + "@" + config.Host + ":" + string(config.Port) + "?database=" + config.Database)
	case "mysql", "mariadb":
		db, err = gorm.Open("mysql", config.Username + ":" + config.Password + "@tcp(" + config.Host + ")/" + config.Database + "?charset=utf8&parseTime=True&loc=Local")
		if err == nil {
			db.Set("gorm:table_options", "ENGINE=InnoDB")
		}
	case "postgre", "postgres", "postgresql":
		db, err = gorm.Open("postgres", "host=" + config.Host + " port=" + string(config.Port) + " user=" + config.Username + " dbname=" + config.Database + " password=" + config.Password)
	case "sqlite", "sqlite3":
		db, err = gorm.Open("sqlite3", config.Database)
	}
	if err != nil {
		panic("Failed to connect database")
	}
}

func Close() (err error) {
	return db.Close()
}
