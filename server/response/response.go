package response

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func Err(c *fiber.Ctx, status int, msg string, errs error) error {
	return c.Status(status).JSON(fiber.Map{
		"message": msg,
		"errors":  errs.Error(),
	})
}

func ErrNotFound(c *fiber.Ctx, err error) error {
	return Err(c, http.StatusNotFound, "Not Found", err)
}

func ErrInternalServer(c *fiber.Ctx, err error) error {
	return Err(c, http.StatusInternalServerError, "Internal Server Error", err)
}

func ErrBadRequest(c *fiber.Ctx, err error) error {
	return Err(c, http.StatusBadRequest, "Bad Request", err)
}

func ErrUnprocessableEntity(c *fiber.Ctx, errs map[string]string) error {
	return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
		"message": "Unprocessable Entity",
		"errors":  errs,
	})
}

func OK[T any](c *fiber.Ctx, data T) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "OK",
		"data":    data,
	})
}
