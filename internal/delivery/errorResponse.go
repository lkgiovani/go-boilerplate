package delivery

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
)

func NewErrorHandler() func(c *fiber.Ctx, err error) error {
	return Error
}

var codes = map[string]int{
	errors.ECONFLICT:       http.StatusConflict,
	errors.EINVALID:        http.StatusBadRequest,
	errors.ENOTFOUND:       http.StatusNotFound,
	errors.ENOTIMPLEMENTED: http.StatusNotImplemented,
	errors.EUNAUTHORIZED:   http.StatusUnauthorized,
	errors.EINTERNAL:       http.StatusInternalServerError,
	errors.EDUPLICATION:    http.StatusConflict,
	errors.EBADREQUEST:     http.StatusBadRequest,
	errors.EFORBIDDEN:      http.StatusForbidden,
	errors.ETIMEOUT:        http.StatusRequestTimeout,
	errors.EUNAVAILABLE:    http.StatusServiceUnavailable,
}

func Error(c *fiber.Ctx, err error) error {
	code, message := errors.ErrorCode(err), errors.ErrorMessage(err)
	if code == errors.EINTERNAL {
		LogError(c, err)
	}

	return c.Status(ErrorStatusCode(code)).JSON(fiber.Map{
		"status":  "error",
		"message": message,
	})
}

func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

func LogError(c *fiber.Ctx, err error) {
	log.Printf("[http] error: %s %s: %s ", c.Method(), c.Path(), err)
}
