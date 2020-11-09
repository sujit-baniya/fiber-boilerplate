package controllers

import (
	"github.com/gofiber/fiber/v2"
	. "github.com/sujit-baniya/fiber-boilerplate/app"
	"github.com/sujit-baniya/fiber-boilerplate/auth"
	"github.com/sujit-baniya/fiber-boilerplate/libraries"
	"github.com/sujit-baniya/fiber-boilerplate/models"
	"strconv"
)

func PlaceOrderFromPaypal(c *fiber.Ctx) error {
	var order models.Payment
	_ = c.BodyParser(&order)
	if amount, err := strconv.ParseFloat(order.Amount, 32); err != nil || order.Amount == "" || amount < 5 {
		return c.Status(400).JSON(fiber.Map{
			"validationError": true,
		})
	}
	pm, err := models.GetPaymentMethodBySlug("paypal")
	if err != nil {
		Flash.WithError(c, fiber.Map{
			"message": "Cannot make order at the moment",
		})
		return c.Redirect("/")

	}
	user, _ := auth.User(c)
	if user.EmailVerified != true {
		return c.Status(400).JSON(fiber.Map{
			"validationError": true,
			"message":         "Email is not verified, ",
		})
	}
	order.PaymentMethodID = pm.ID
	order.Currency = pm.Currency
	order.UserID = user.ID
	order.Status = "PENDING"
	_, _ = order.Create()
	err = libraries.CreateOrder(&order, user)
	DB.Save(&order)

	if err != nil {
		Flash.WithError(c, fiber.Map{
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

func PostOrderResponseFromPaypal(c *fiber.Ctx) {

}

func ValidateOrderFromPaypal(c *fiber.Ctx) error {
	amount := c.Params("amount")

	if _, err := strconv.ParseFloat(amount, 32); err != nil || amount == "" {
		return c.Status(400).JSON(fiber.Map{
			"validationError": true,
		})
	}
	return c.JSON(fiber.Map{
		"validationError": false,
	})
}

func GetOrderDetailFromPaypal(c *fiber.Ctx) error {
	order, err := libraries.GetOrder(c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}
	return c.JSON(order)
}

func PostOrderCancelResponseFromPaypal(c *fiber.Ctx) error {
	order, err := libraries.GetOrder(c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}
	p, err := models.GetPaymentByGatewayOrderID(order.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
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
	order, err := libraries.GetOrder(c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}
	p, err := models.GetPaymentByGatewayOrderID(order.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": err,
		})
	}
	p.PayPalOrderDetail = order
	p.GatewayOrderID = order.ID
	p.UpdatePaymentStatusByGatewayOrderID("APPROVED")
	user, _ := models.GetUserById(p.UserID)
	user.AddAmount(p.Amount)
	return c.JSON(p)
}
