package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dwadp/attendance-api/internal"
	holidayTypes "github.com/dwadp/attendance-api/internal/holiday/types"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/server/response"
	"github.com/dwadp/attendance-api/server/validator"
	"github.com/dwadp/attendance-api/store"
	"github.com/gofiber/fiber/v2"
)

func handleGetListHoliday(store store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		holidays, err := store.FindAllHoliday(c.UserContext())
		if err != nil {
			return response.ErrInternalServer(c, fmt.Errorf("failed to get list of holidays: %v", err))
		}
		return response.OK(c, holidays)
	}
}

func handleCreateHoliday(store store.Store, v *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var create models.UpsertHoliday

		if err := c.BodyParser(&create); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse request body: %v", err))
		}

		if create.Type == holidayTypes.NationalHoliday && create.Weekday != holidayTypes.None {
			return response.ErrBadRequest(c, internal.ErrNationalHolidayShouldNotHaveWeekday)
		}

		if create.Type == holidayTypes.Weekend && create.Date.Valid {
			return response.ErrBadRequest(c, internal.ErrWeekendShouldNotHaveDate)
		}

		if create.Date.Valid || create.Weekday != holidayTypes.None {
			existing, err := store.FindHolidayByDateOrWeekday(c.UserContext(), create.Date.T, create.Weekday)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return response.ErrInternalServer(c, err)
				}
			}

			if existing != nil {
				return response.OK(c, existing)
			}
		}

		if err := v.Validate(create); err != nil {
			return response.ErrUnprocessableEntity(c, v.SerializeErrors(err, create))
		}

		holiday, err := store.CreateHoliday(c.UserContext(), create)
		if err != nil {
			if errors.Is(err, internal.ErrWeekendShouldNotHaveDate) || errors.Is(err, internal.ErrNationalHolidayShouldNotHaveWeekday) {
				return response.ErrBadRequest(c, err)
			}

			return response.ErrInternalServer(c, fmt.Errorf("unable to create holiday: %v", err))
		}

		return response.OK(c, holiday)
	}
}

func handleGetDetailHoliday(store store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		holiday, err := store.FindHolidayByID(c.UserContext(), uint(id))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return response.ErrNotFound(c, fmt.Errorf("holiday ID %d could not be found", id))
			}

			return response.ErrInternalServer(c, fmt.Errorf("failed to find holiday ID: %q", err.Error()))
		}

		return response.OK(c, holiday)
	}
}

func handleDeleteHoliday(store store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse holiday ID: %v", err))
		}

		if err := store.DeleteHoliday(c.UserContext(), uint(id)); err != nil {
			return response.ErrInternalServer(c, fmt.Errorf("failed to delete holiday ID %d: %v", id, err))
		}

		return response.OK[any](c, nil)
	}
}

func handleUpdateHoliday(store store.Store, v *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse holiday ID: %v", err))
		}

		var update models.UpsertHoliday
		if err := c.BodyParser(&update); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse request body: %v", err))
		}

		if update.Type == holidayTypes.NationalHoliday && update.Weekday != holidayTypes.None {
			return response.ErrBadRequest(c, internal.ErrNationalHolidayShouldNotHaveWeekday)
		}

		if update.Type == holidayTypes.Weekend && update.Date.Valid {
			return response.ErrBadRequest(c, internal.ErrWeekendShouldNotHaveDate)
		}

		err = v.Validate(update)
		if err != nil {
			return response.ErrUnprocessableEntity(c, v.SerializeErrors(err, update))
		}

		holiday, err := store.UpdateHoliday(c.UserContext(), uint(id), update)
		if err != nil {
			if errors.Is(err, internal.ErrWeekendShouldNotHaveDate) || errors.Is(err, internal.ErrNationalHolidayShouldNotHaveWeekday) {
				return response.ErrBadRequest(c, err)
			}

			return response.ErrInternalServer(c, fmt.Errorf("failed to update holiday ID %d: %v", id, err))
		}

		return response.OK(c, holiday)
	}
}

func RegisterHolidayHandlers(router fiber.Router, store store.Store, v *validator.Validator) {
	router.Get("/", handleGetListHoliday(store))
	router.Post("/", handleCreateHoliday(store, v))
	router.Get("/:id", handleGetDetailHoliday(store))
	router.Delete("/:id", handleDeleteHoliday(store))
	router.Put("/:id", handleUpdateHoliday(store, v))
}
