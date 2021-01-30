package models

import (
	"strings"

	"github.com/sujit-baniya/fiber-boilerplate/app"
	"gorm.io/gorm"
)

type PaymentMethod struct {
	*gorm.Model
	Name     string `json:"name" gorm:"name"`           //nolint:gofmt
	Slug     string `json:"slug" gorm:"slug"`           //nolint:gofmt
	IsActive bool   `json:"is_active" gorm:"is_active"` //nolint:gofmt
	Currency string `json:"currency" gorm:"currency"`
}

func GetPaymentMethodBySlug(slug string) (*PaymentMethod, error) {
	var pm PaymentMethod
	if err := app.Http.Database.DB.Where(&PaymentMethod{Slug: slug}).FirstOrCreate(&pm).Error; err != nil {
		return nil, err
	}
	return &pm, nil
}

func (u *PaymentMethod) BeforeCreate(db *gorm.DB) (err error) {
	if u.Name == "" {
		u.Name = strings.ToTitle(u.Slug)
	}
	if u.Currency == "" {
		u.Currency = "USD"
	}
	return
}
