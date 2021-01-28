package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/verify-rest/app"
	"github.com/sujit-baniya/verify-rest/pkg/auth"
	"github.com/sujit-baniya/verify-rest/pkg/models"
)

func PlaceOrderFromPaypal(c *fiber.Ctx) error {
	paypal := &models.PayPal{&app.Http.PayPal}
	var order models.Payment
	c.BodyParser(&order)
	if amount, err := strconv.ParseFloat(order.Amount, 32); err != nil || order.Amount == "" || amount < 5 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"validationError": true,
		})
	}
	pm, err := models.GetPaymentMethodBySlug("paypal")
	if err != nil {
		fmt.Println(err.Error())
		app.Http.Flash.WithError(c, fiber.Map{
			"message": "Cannot make order at the moment",
		})
		return c.Redirect("/")
	}
	user, _ := auth.User(c)
	if user.EmailVerified != true {
		c.Status(400)
		return c.JSON(fiber.Map{
			"validationError": true,
			"message":         "Email is not verified, ",
		})
	}
	order.PaymentMethodID = pm.ID
	order.Currency = pm.Currency
	order.UserID = user.ID
	order.Status = "PENDING"
	order.Create()
	err = paypal.CreateOrder(&order, user)
	app.Http.Database.DB.Save(&order)

	if err != nil {
		fmt.Println(err.Error())
		app.Http.Flash.WithError(c, fiber.Map{
			"message": "Cannot make order at the moment",
		})
		return c.Redirect("/")
	}
	return c.JSON(fiber.Map{
		"ack": true,
		"data": fiber.Map{
			"id": order.PayPalOrderDetail.ID,
		},
	})
}

func PostOrderResponseFromPaypal(c *fiber.Ctx) error {
	return nil
}

func ValidateOrderFromPaypal(c *fiber.Ctx) error {
	amount := c.Params("amount")

	if _, err := strconv.ParseFloat(amount, 32); err != nil || amount == "" {
		c.Status(400)
		return c.JSON(fiber.Map{
			"validationError": true,
		})
	}
	return c.JSON(fiber.Map{
		"validationError": false,
	})
}

func GetOrderDetailFromPaypal(c *fiber.Ctx) error {
	paypal := &models.PayPal{&app.Http.PayPal}
	order, err := paypal.GetOrder(c.Params("id"))
	if err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}
	return c.JSON(order)
}

func PostOrderCancelResponseFromPaypal(c *fiber.Ctx) error {
	paypal := &models.PayPal{&app.Http.PayPal}
	order, err := paypal.GetOrder(c.Params("id"))
	if err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}
	p, err := models.GetPaymentByGatewayOrderID(order.ID)
	if err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}
	p.PayPalOrderDetail = order
	p.GatewayOrderID = order.ID
	p.UpdatePaymentStatusByGatewayOrderID("CANCELED")
	return c.JSON(order)
}

func PostOrderSuccessResponseFromPaypal(c *fiber.Ctx) error {
	paypal := &models.PayPal{&app.Http.PayPal}
	order, err := paypal.GetOrder(c.Params("id"))
	if err != nil {
		c.Status(500)
		return c.JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}
	p, err := models.GetPaymentByGatewayOrderID(order.ID)
	if err != nil {
		return err
	}
	p.PayPalOrderDetail = order
	p.GatewayOrderID = order.ID
	user, _ := models.GetUserById(p.UserID)
	user.AddAmount(p.Amount, p.BalanceAdded)
	p.UpdatePaymentStatusByGatewayOrderID(order.Status)
	return c.JSON(p)
}
