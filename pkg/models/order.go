package models

import "github.com/plutov/paypal/v3"

type Order struct {
	InvoiceID         string
	Amount            string
	Currency          string
	PayPalOrderDetail *paypal.Order
	OrderUrl          string
}
