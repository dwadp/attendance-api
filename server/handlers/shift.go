package handlers

import (
	"errors"
	"fmt"
	"github.com/dwadp/attendance-api/internal/holiday"
	shiftInt "github.com/dwadp/attendance-api/internal/shift"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/server/response"
	"github.com/dwadp/attendance-api/server/validator"
	"github.com/dwadp/attendance-api/store"
	"github.com/dwadp/attendance-api/store/db"
	"github.com/gofiber/fiber/v2"
)

func handleGetListShifts(store store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		shifts, err := store.FindAllShifts(c.UserContext())
		if err != nil {
			return response.ErrInternalServer(c, fmt.Errorf("failed to get list of shifts: %v", err))
		}
		return response.OK(c, shifts)
	}
}

func handleCreateShift(store store.Store, v *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var create models.UpsertShift

		if err := c.BodyParser(&create); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse request body: %w", err))
		}

		err := v.Validate(create)
		if err != nil {
			return response.ErrUnprocessableEntity(c, v.SerializeErrors(err, create))
		}

		shift, err := store.CreateShift(c.UserContext(), create)
		if err != nil {
			return response.ErrInternalServer(c, fmt.Errorf("unable to create shift: %w", err))
		}

		return response.OK(c, shift)
	}
}

func handleGetDetailShift(store store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		shift, err := store.FindShiftByID(c.UserContext(), uint(id))
		if err != nil {
			var errDataNotFound *db.ErrDataNotFound[uint]

			if errors.As(err, &errDataNotFound) {
				return response.ErrNotFound(c, fmt.Errorf("shift ID %d could not be found", id))
			}

			return response.ErrInternalServer(c, fmt.Errorf("failed to find shift ID: %q", err.Error()))
		}

		return response.OK(c, shift)
	}
}

func handleDeleteShift(store store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse shift ID: %v", err))
		}

		if err := store.DeleteShift(c.UserContext(), uint(id)); err != nil {
			return response.ErrInternalServer(c, fmt.Errorf("failed to delete shift ID %d: %v", id, err))
		}

		return response.OK[any](c, nil)
	}
}

func handleUpdateShift(store store.Store, v *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse shift ID: %v", err))
		}

		var update models.UpsertShift
		if err := c.BodyParser(&update); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse request body: %v", err))
		}

		err = v.Validate(update)
		if err != nil {
			return response.ErrUnprocessableEntity(c, v.SerializeErrors(err, update))
		}

		shift, err := store.UpdateShift(c.UserContext(), uint(id), update)
		if err != nil {
			return response.ErrInternalServer(c, fmt.Errorf("failed to update shift ID %d: %v", id, err))
		}

		return response.OK(c, shift)
	}
}

func handleEmployeeShiftAssignment(store store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var assign models.AssignEmployeeShift
		if err := c.BodyParser(&assign); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse request body: %v", err))
		}

		holidays, err := store.FindAllHolidays(c.UserContext(), assign.Date.T)
		if err != nil {
			return response.ErrInternalServer(c, err)
		}

		hService := holiday.NewService(holidays)
		s := shiftInt.NewService(store, hService)

		result, err := s.AssignEmployee(c.UserContext(), assign)
		if err != nil {
			if errors.Is(err, holiday.ErrIsOnHoliday) {
				return response.ErrBadRequest(c, err)
			}
			return response.ErrInternalServer(c, err)
		}

		return response.OK(c, result)
	}
}

func RegisterShift(router fiber.Router, store store.Store, v *validator.Validator) {
	router.Get("/", handleGetListShifts(store))
	router.Post("/", handleCreateShift(store, v))
	router.Get("/:id", handleGetDetailShift(store))
	router.Delete("/:id", handleDeleteShift(store))
	router.Put("/:id", handleUpdateShift(store, v))
	router.Post("/assign", handleEmployeeShiftAssignment(store))
}
