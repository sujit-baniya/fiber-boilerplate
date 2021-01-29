package migrations

import (
	"log"

	"github.com/sujit-baniya/fiber-boilerplate/app"
	"github.com/sujit-baniya/fiber-boilerplate/pkg/models"
)

func Migrate() {
	log.Println("Initiating migration...")
	err := app.Http.Database.DB.Migrator().AutoMigrate(
		&models.User{},
		&models.UserMeta{},
		&models.UserSetting{},
		&models.File{},
		&models.PaymentMethod{},
		&models.Payment{},
		&models.UserFile{},
		&models.Transaction{},
		&models.UserTransactionLog{},
	)
	if err != nil {
		panic(err)
	}
	log.Println("Migration Completed...")
}
