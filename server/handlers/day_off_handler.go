package handlers

import (
	"errors"
	"fmt"
	"github.com/dwadp/attendance-api/internal"
	"github.com/dwadp/attendance-api/internal/dayoff"
	"github.com/dwadp/attendance-api/internal/holiday"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/server/response"
	"github.com/dwadp/attendance-api/server/validator"
	"github.com/dwadp/attendance-api/store"
	"github.com/gofiber/fiber/v2"
)

func handleRequestDayOff(v *validator.Validator, service *dayoff.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request models.DayOffRequest
		if err := c.BodyParser(&request); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse request body: %v", err))
		}

		if err := v.Validate(request); err != nil {
			return response.ErrUnprocessableEntity(c, v.SerializeErrors(err, request))
		}

		result, err := service.Create(c.UserContext(), request)
		if err != nil {
			if errors.Is(err, internal.ErrShiftExists) {
				return response.ErrBadRequest(c, err)
			} else if errors.Is(err, internal.ErrIsOnHoliday) {
				return response.ErrBadRequest(c, err)
			}

			return response.ErrInternalServer(c, err)
		}

		return response.OK(c, result)
	}
}

func RegisterDayOff(router fiber.Router, store store.Store, v *validator.Validator) {
	holidayService := holiday.NewService(store)
	dayOffService := dayoff.NewService(store, holidayService)

	router.Post("/", handleRequestDayOff(v, dayOffService))
}
