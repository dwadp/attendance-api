package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dwadp/attendance-api/internal"
	"github.com/dwadp/attendance-api/internal/holiday"
	"github.com/dwadp/attendance-api/internal/shift"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/server/response"
	"github.com/dwadp/attendance-api/server/validator"
	"github.com/dwadp/attendance-api/store"
	"github.com/gofiber/fiber/v2"
)

func handleEmployeeShiftAssignment(service *shift.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var assign models.AssignEmployeeShift
		if err := c.BodyParser(&assign); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse request body: %v", err))
		}

		result, err := service.AssignEmployee(c.UserContext(), assign)
		if err != nil {
			if errors.Is(err, internal.ErrIsOnHoliday) {
				return response.ErrBadRequest(c, err)
			} else if errors.Is(err, internal.ErrDayOffExistsOnDate) {
				return response.ErrBadRequest(c, err)
			}

			return response.ErrInternalServer(c, err)
		}

		return response.OK(c, result)
	}
}

func handleEmployeeShiftUnassignment(store store.Store, v *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var unassign models.UnassignEmployeeShift
		if err := c.BodyParser(&unassign); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse request body: %v", err))
		}

		if err := v.Validate(unassign); err != nil {
			return response.ErrUnprocessableEntity(c, v.SerializeErrors(err, unassign))
		}

		if err := store.DeleteEmployeeShift(c.UserContext(), unassign); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return response.ErrInternalServer(c, fmt.Errorf("failed to delete employee shift: %v", err))
			}
		}

		return response.OK[any](c, nil)
	}
}

func RegisterEmployeeShift(router fiber.Router, store store.Store, v *validator.Validator) {
	holidayService := holiday.NewService(store)
	shiftService := shift.NewService(store, holidayService)

	router.Post("/assign", handleEmployeeShiftAssignment(shiftService))
	router.Post("/unassign", handleEmployeeShiftUnassignment(store, v))
}
