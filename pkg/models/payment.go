package models

import (
	"github.com/plutov/paypal/v3"
	"github.com/sujit-baniya/fiber-boilerplate/app"
	"gorm.io/gorm"
)

type Payment struct {
	*gorm.Model
	PaymentMethodID    uint          `json:"payment_method_id" gorm:"payment_method_id"` //nolint:gofmt
	PaymentMethod      PaymentMethod `gorm:"foreignkey:PaymentMethodID"`
	UserID             uint          `json:"user_id" gorm:"user_id"` //nolint:gofmt
	User               User          `gorm:"foreignKey:UserID"`
	Amount             string        `json:"amount" gorm:"amount"`
	Status             string        `json:"status" gorm:"status"`
	GatewayOrderID     string        `json:"gateway_order_id" gorm:"gateway_order_id"`
	GatewayOrderStatus string        `json:"gateway_order_status" gorm:"gateway_order_status"`
	Currency           string        `json:"currency" gorm:"currency"`
	BalanceAdded       bool          `json:"balance_added" gorm:"balance_added"`
	PayPalOrderDetail  *paypal.Order `gorm:"-"`
}

func (l *Payment) Create() (*Payment, error) {
	app.Http.Database.DB.Create(l)
	return l, nil
}

func GetPaymentByGatewayOrderID(id string) (*Payment, error) {
	var p Payment
	if err := app.Http.Database.DB.Preload("User").Preload("PaymentMethod").Where(&Payment{GatewayOrderID: id}).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil //nolint:wsl
}

func (l *Payment) UpdatePaymentStatusByGatewayOrderID(status string) {
	l.Status = status
	l.GatewayOrderStatus = l.PayPalOrderDetail.Status
	l.BalanceAdded = true
	app.Http.Database.DB.Save(&l)
}
