package error_handler

import (
	"fmt"
	"github.com/gofiber/fiber"
	"io"
	"os"
	"strconv"
	"strings"
)

// Config ...
type Config struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool
	// Custom error handler
	// Optional. Default: nil
	Handler func(*fiber.Ctx, error, func(...interface{}))
	// Log all errors to output
	// Optional. Default: false
	Log bool
	// Output is a writer where logs are written
	// Default: os.Stderr
	Output io.Writer
	// Use c.Render for content-type html
	// Optional. Default: false
	UseTemplate bool
}

// Send error message as JSON
func handleJSON(c *fiber.Ctx, args ...interface{}) {
	l := len(args)
	var httpErr HTTPError

	httpErr = NewHttpError(fiber.StatusInternalServerError, "Internal Server Error", nil)
	if l > 0 {
		if he, ok := args[0].(HTTPError); ok {
			httpErr = he
		} else if e, ok := args[0].(error); ok {
			httpErr = NewHttpError(fiber.StatusInternalServerError, e.Error(), e.Error())
		} else if s, ok := args[0].(string); ok {
			httpErr = NewHttpError(fiber.StatusInternalServerError, s, s)
		} else {
			httpErr = NewHttpError(fiber.StatusInternalServerError, "Internal Server Error", args[0])
		}
	}

	c.Status(httpErr.StatusCode())

	if httpErr.Data() != nil {
		c.JSON(fiber.Map{
			"message": httpErr.Message(),
			"error":   httpErr.Data(),
		})
	} else {
		c.JSON(fiber.Map{
			"message": httpErr.Message(),
		})
	}
}

// Render template based on args
// Posible args combinations are:
// handleTemplate(*fiber.Ctx, string)
// handleTemplate(*fiber.Ctx, string, error)
// handleTemplate(*fiber.Ctx, error)
func handleTemplate(c *fiber.Ctx, args ...interface{}) {
	l := len(args)
	var httpErr HTTPError = NewHttpError(fiber.StatusInternalServerError, "Internal Server Error", nil)
	var view = "500"

	if l > 0 {
		if s, ok := args[0].(string); ok {
			view = s

			if l >= 2 {
				if he, ok := args[1].(HTTPError); ok {
					httpErr = he
				} else if e, ok := args[1].(error); ok {
					httpErr = NewHttpError(fiber.StatusInternalServerError, e.Error(), e)
				} else if s, ok := args[1].(string); ok {
					httpErr = NewHttpError(fiber.StatusInternalServerError, s, s)
				} else {
					httpErr = NewHttpError(fiber.StatusInternalServerError, "Internal Server Error", args[1])
				}
			}
		} else if he, ok := args[0].(HTTPError); ok {
			view = strconv.Itoa(he.StatusCode())
			httpErr = he
		} else if e, ok := args[0].(error); ok {
			httpErr = NewHttpError(fiber.StatusInternalServerError, e.Error(), e)
			view = strconv.Itoa(httpErr.StatusCode())
		}
	}
	view = fmt.Sprintf("errors/%s", view)
	c.Status(httpErr.StatusCode()).Render(view, fiber.Map{
		"error": httpErr,
	})
}

// Send error message as plain text
func handlePlainText(c *fiber.Ctx, args ...interface{}) {
	l := len(args)

	if l > 0 {
		if s, ok := args[0].(string); ok {
			c.Status(fiber.StatusInternalServerError).SendString(s)
			return
		} else if he, ok := args[0].(HTTPError); ok {
			c.Status(he.StatusCode()).SendString(he.Message())
			return
		} else if e, ok := args[0].(error); ok {
			c.Status(fiber.StatusInternalServerError).SendString(e.Error())
			return
		}
	}

	c.SendStatus(fiber.StatusInternalServerError)
}

// Decide the content type to be used based on `Content-Type` or `Accept` header.
// Only accept `text/plain`, `*/json`, `*/html` or `*/xhtml-xml`
// If `Accept` header contains multiple values, the first to come up will be used with `text/plain` as fallback value.
func getPreferedContentType(c *fiber.Ctx) (ct string) {
	// default text/plain
	ct = fiber.MIMETextPlain
	header := strings.ToLower(c.Get(fiber.HeaderContentType))
	if header != "" {
		if factorSign := strings.IndexByte(header, ';'); factorSign != -1 {
			ct = header[:factorSign]
			return
		}
		ct = header
	} else {
		// if 'Content-Type' is empty then use 'Accept'
		header = strings.ToLower(c.Get(fiber.HeaderAccept))
		if factorSign := strings.IndexByte(header, ';'); factorSign != -1 {
			header = header[:factorSign]
		}
		vals := strings.Split(header, ",")

		for _, val := range vals {
			if val == fiber.MIMETextPlain {
				return
			} else if strings.HasSuffix(val, "json") {
				ct = fiber.MIMEApplicationJSON
				return
			} else if strings.HasSuffix(val, "html") || strings.HasSuffix(val, "xhtml+xml") {
				ct = fiber.MIMETextHTML
				return
			}
		}
	}
	return
}

// New ...
func New(config ...Config) func(*fiber.Ctx) {
	// Init config
	var cfg Config
	// Set config if provided
	if len(config) > 0 {
		cfg = config[0]
	}
	if cfg.Output == nil {
		cfg.Output = os.Stderr
	}

	// Return middleware handler
	return func(c *fiber.Ctx) {
		// default handler
		errHandler := func(args ...interface{}) {
			ct := getPreferedContentType(c)

			if ct == fiber.MIMEApplicationJSON {
				handleJSON(c, args...)
				return
			} else if ct == fiber.MIMETextHTML {
				// use template
				if cfg.UseTemplate {
					handleTemplate(c, args...)
				} else {
					// use json if template is not used
					handleJSON(c, args...)
				}
				return
			} else {
				handlePlainText(c, args...)
				return
			}
		}

		// Filter request to skip middleware
		if cfg.Filter != nil && cfg.Filter(c) {
			c.Next()
			return
		}
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				// Log error
				if cfg.Log {
					cfg.Output.Write([]byte(err.Error() + "\n"))
				}
				if cfg.Handler != nil {
					cfg.Handler(c, err, errHandler)
				} else {
					errHandler(err)
				}
			}
		}()
		c.Next()
		if c.Error() != nil {
			if cfg.Log {
				cfg.Output.Write([]byte(c.Error().Error() + "\n"))
			}

			if cfg.Handler != nil {
				cfg.Handler(c, c.Error(), errHandler)
			} else {
				errHandler(c.Error())
			}
		}
	}
}
