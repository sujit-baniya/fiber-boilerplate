package models

import (
	"errors"
	. "github.com/itsursujit/fiber-boilerplate/app"
	"github.com/jinzhu/gorm"
)

type RegisterForm struct {
	gorm.Model
	Name      string `json:"name" gorm:"name" form:"name" `
	Email     string `json:"email" gorm:"email" form:"email" validate:"required|email"`
	Password  string `json:"-" gorm:"password" form:"password" validate:"required"`
	CPassword string `json:"-" form:"c_password" validate:"required|eq_field:password" gorm:"-"`
}

func (RegisterForm) TableName() string {
	return "users"
}

func (l *RegisterForm) Signup() (*RegisterForm, error) {

	user, _ := GetUserByEmail(l.Email)
	if user != nil {
		return nil, errors.New("User Already Exists")
	}
	l.Password, _ = Hash.Create(l.Password)
	DB.Create(l)
	return l, nil
}

func (l *RegisterForm) ResetPassword() (*RegisterForm, error) {
	l.Password, _ = Hash.Create(l.Password)
	DB.Updates(&l)
	return l, nil
}
