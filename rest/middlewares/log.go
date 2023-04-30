package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oarkflow/log"
)

func NewLog(logger *log.Logger, skip func(c *fiber.Ctx) bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		next := c.Next()

		if skip != nil && skip(c) {
			return nil
		}

		end := time.Now()
		latency := end.Sub(start)

		status := c.Response().StatusCode()
		msg := "Request"
		if next != nil {
			msg = next.Error()
		}

		var e *log.Entry
		switch {
		case status >= 400 && status < 500:
			e = logger.Warn()
		case status >= 500:
			e = logger.Error()
		default:
			e = logger.Info()
		}
		e.Int("status", status).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("ip", c.IP()).
			Dur("latency", latency).
			Str("user_agent", c.Get(fiber.HeaderUserAgent)).
			Msg(msg)

		return nil
	}
}
