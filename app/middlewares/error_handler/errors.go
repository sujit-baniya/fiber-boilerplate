package error_handler

import (
	"fmt"
	"github.com/gofiber/fiber"
)

type HTTPError interface {
	// HTTP Status code to respond with.
	// Optional. Default: 500
	StatusCode() int
	// Error message
	// Optional. Default: ""
	Message() string
	// Data to respond with or to bind to views.
	// Optional. Default: nil
	Data() interface{}
	Error() string
}

type httpError struct {
	statusCode int
	message    string
	data       interface{}
}

func NewHttpError(statusCode int, message string, data interface{}) *httpError {
	if statusCode == 0 {
		statusCode = fiber.StatusInternalServerError
	}
	return &httpError{
		statusCode: statusCode,
		message:    message,
		data:       data,
	}
}

func (he *httpError) StatusCode() int {
	return he.statusCode
}

func (he *httpError) Message() string {
	return he.message
}

func (he *httpError) Data() interface{} {
	return he.data
}

func (he *httpError) Error() string {
	return fmt.Sprintf("statusCode: %d, message: %s", he.statusCode, he.message)
}
