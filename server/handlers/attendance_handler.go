package handlers

import (
	"fmt"
	"github.com/dwadp/attendance-api/internal/attendance"
	"github.com/dwadp/attendance-api/internal/attendance/types"
	"github.com/dwadp/attendance-api/internal/holiday"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/server/response"
	"github.com/dwadp/attendance-api/server/validator"
	"github.com/dwadp/attendance-api/store"
	"github.com/gofiber/fiber/v2"
)

func handleRequestClockIn(service *attendance.Service, v *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request models.AttendanceRequest
		if err := c.BodyParser(&request); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse request: %w", err))
		}

		if err := v.Validate(&request); err != nil {
			return response.ErrUnprocessableEntity(c, v.SerializeErrors(err, request))
		}

		request.Type = types.ClockIn

		result, err := service.ClockIn(c.UserContext(), request)
		if err != nil {
			return response.ErrInternalServer(c, err)
		}

		return response.OK(c, result)
	}
}

func handleRequestClockOut(service *attendance.Service, v *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request models.AttendanceRequest
		if err := c.BodyParser(&request); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse request: %w", err))
		}

		if err := v.Validate(&request); err != nil {
			return response.ErrUnprocessableEntity(c, v.SerializeErrors(err, request))
		}

		request.Type = types.ClockIn

		result, err := service.ClockOut(c.UserContext(), request)
		if err != nil {
			return response.ErrInternalServer(c, err)
		}

		return response.OK(c, result)
	}
}

func RegisterAttendance(router fiber.Router, store store.Store, v *validator.Validator) {
	holidayService := holiday.NewService(store)
	attendanceService := attendance.NewService(store, holidayService)

	router.Post("/clock-in", handleRequestClockIn(attendanceService, v))
	router.Post("/clock-out", handleRequestClockOut(attendanceService, v))
}