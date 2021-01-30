package models

import (
	"errors"

	"github.com/sujit-baniya/fiber-boilerplate/app"
	"gorm.io/gorm"
)

type RegisterForm struct {
	gorm.Model
	FirstName string `json:"first_name" gorm:"first_name" form:"first_name" `
	LastName  string `json:"last_name" gorm:"last_name" form:"last_name" `
	Email     string `json:"email" gorm:"email" form:"email" validate:"required|email"`
	Password  string `json:"password" gorm:"password" form:"password" validate:"required"`
	CPassword string `json:"c_password" form:"c_password" validate:"required|eq_field:password" gorm:"-"`
}

func (RegisterForm) TableName() string {
	return "users"
}

func (l *RegisterForm) Signup() (*RegisterForm, error) {

	user, _ := GetUserByEmail(l.Email)
	if user != nil {
		return nil, errors.New("User Already Exists")
	}
	l.Password, _ = app.Http.Hash.Create(l.Password)
	app.Http.Database.DB.Create(l)
	return l, nil
}

func (l *RegisterForm) ResetPassword() (*RegisterForm, error) {
	l.Password, _ = app.Http.Hash.Create(l.Password)
	app.Http.Database.DB.Updates(&l)
	return l, nil
}
