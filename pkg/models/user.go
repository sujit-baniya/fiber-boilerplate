package models

import (
	"strconv"
	"time"

	"github.com/sujit-baniya/fiber-boilerplate/app"
	"gorm.io/gorm"
)

type User struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	FirstName     string         `json:"first_name" gorm:"first_name"` //nolint:gofmt
	LastName      string         `json:"last_name" gorm:"last_name"`   //nolint:gofmt
	Email         string         `json:"email" gorm:"email"`
	Password      string         `json:"-" gorm:"password"`
	Balance       float32        `json:"balance" gorm:"balance"`
	EmailVerified bool           `json:"email_verified" gorm:"email_verified"`
	Currency      string         `json:"currency" gorm:"currency"`
	IsAdmin       bool           `json:"is_admin" gorm:"is_admin"`
	Files         []File         `gorm:"many2many:user_files;" json:"files,omitempty"`
	Metas         []UserMeta     `gorm:"many2many:user_meta;" json:"metas,omitempty"`
	UserSetting   UserSetting    `gorm:"-" json:"settings,omitempty"`
}

type Role struct {
	ID          uint   `gorm:"primarykey"`
	Name        string `json:"name" gorm:"name"`               //nolint:gofmt
	Slug        string `json:"slug" gorm:"slug"`               //nolint:gofmt
	Description string `json:"description" gorm:"description"` //nolint:gofmt
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func AllUsers() []User {
	var users []User
	app.Http.Database.Find(&users)
	return users
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := app.Http.Database.Preload("Metas").Where(&User{Email: email}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil //nolint:wsl
}

func GetVerifiedUserByEmail(email string) (*User, error) {
	var user User
	if err := app.Http.Database.Preload("Metas").Where(&User{Email: email, EmailVerified: true}).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil //nolint:wsl
}

func GetUserById(id interface{}) (*User, error) {
	var user User
	if err := app.Http.Database.Preload("Files").Where("id = ? ", id).First(&user).Error; err != nil {
		return nil, err
	}
	settings, err := user.Settings()
	if err != nil {
		// do something sensible
	}
	user.UserSetting = settings
	return &user, nil
}

func (u *User) Update() error {
	if u.ID != 0 {
		if err := app.Http.Database.Updates(&u).Error; err != nil {
			return err
		}
	} else {
		if err := app.Http.Database.Where(&User{Email: u.Email}).Updates(&u).Error; err != nil {
			return err
		}
	}
	app.Http.Database.First(&u)
	u.Settings()
	return nil
}

func (u *User) AddAmount(amount string, AlreadyAdded bool) {
	if AlreadyAdded {
		return
	}
	value, err := strconv.ParseFloat(amount, 32)
	if err != nil {
		// do something sensible
	}
	u.Balance += float32(value)
	app.Http.Database.Updates(&u)
}

func (u *User) Settings() (UserSetting, error) {
	userSettings := UserSetting{UserID: u.ID}
	err := userSettings.Get()
	if err != nil {
		return UserSetting{}, err
	}
	u.UserSetting = userSettings
	return userSettings, nil
}
