package models

import (
	"fmt"
	"github.com/sujit-baniya/fiber-boilerplate/config"

	"github.com/plutov/paypal/v3"
	"github.com/sujit-baniya/fiber-boilerplate/app"
)

type PayPal struct {
	*config.PayPalConfig
}

func (p *PayPal) GetOrder(id string) (*paypal.Order, error) {
	_, err := p.Client.GetAccessToken()
	if err != nil {
		fmt.Println(err)
	}
	return p.Client.GetOrder(id)
}

func (p *PayPal) CreateOrder(o *Payment, user *User) error {

	_, err := p.Client.GetAccessToken()
	if err != nil {
		fmt.Println(err)
	}
	order, err := p.Client.CreateOrder(
		paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			{
				ReferenceID: fmt.Sprintf("%d", o.ID),
				Amount: &paypal.PurchaseUnitAmount{
					Value:    o.Amount,
					Currency: o.Currency,
				},
			},
		},
		&paypal.CreateOrderPayer{
			Name: &paypal.CreateOrderPayerName{
				GivenName: user.FirstName,
				Surname:   user.LastName,
			},
			EmailAddress: user.Email,
		},
		&paypal.ApplicationContext{
			BrandName: app.Http.Server.Name,
			ReturnURL: fmt.Sprintf("%s/paypal/response", app.Http.Server.Url),
			CancelURL: fmt.Sprintf("%s/paypal/cancel", app.Http.Server.Url),
		},
	)
	if err != nil {
		return err
	}
	p.Client.CaptureOrder(order.ID, paypal.CaptureOrderRequest{})
	o.PayPalOrderDetail = order
	o.GatewayOrderID = order.ID
	o.GatewayOrderStatus = order.Status
	o.Status = "PROCESSING"
	return nil
}
