package models

import (
	"github.com/sujit-baniya/fiber-boilerplate/app"
	"gorm.io/gorm"
)

type UserSetting struct {
	*gorm.Model
	UserID      uint   `gorm:"user_id" json:"user_id"`
	PhoneRegion string `gorm:"phone_region" json:"phone_region"`
	SenderID    string `gorm:"sender_id" json:"sender_id"`
	QueueType   string `gorm:"queue_type" json:"queue_type"` // SHARED | ACCOUNT | CAMPAIGN
}

func (u *UserSetting) UpdateOrCreate() {
	if u.UserID != 0 {
		if app.Http.Database.Where(&UserSetting{UserID: u.UserID}).Updates(&u).RowsAffected == 0 {
			app.Http.Database.Create(&u)
		}
	}
}

func (u *UserSetting) Get() error {
	if err := app.Http.Database.Where(&UserSetting{UserID: u.UserID}).First(&u).Error; err != nil {
		return err
	}
	return nil
}
