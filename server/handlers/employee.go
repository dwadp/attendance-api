package handlers

import (
	"errors"
	"fmt"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/server/response"
	"github.com/dwadp/attendance-api/server/validator"
	"github.com/dwadp/attendance-api/store"
	"github.com/dwadp/attendance-api/store/db"
	"github.com/gofiber/fiber/v2"
)

func handleGetListEmployee(store store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		employees, err := store.FindAllEmployees(c.UserContext())
		if err != nil {
			return response.ErrInternalServer(c, fmt.Errorf("failed to get list of employees: %v", err))
		}
		return response.OK(c, employees)
	}
}

func handleCreateEmployee(store store.Store, v *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var create models.UpsertEmployee

		if err := c.BodyParser(&create); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse request body: %v", err))
		}

		if err := v.Validate(create); err != nil {
			return response.ErrUnprocessableEntity(c, v.SerializeErrors(err, create))
		}

		employee, err := store.CreateEmployee(c.UserContext(), create)
		if err != nil {
			return response.ErrInternalServer(c, fmt.Errorf("unable to create employee: %v", err))
		}

		return c.JSON(fiber.Map{
			"data": employee,
		})
	}
}

func handleGetDetailEmployee(store store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		employee, err := store.FindEmployeeByID(c.UserContext(), uint(id))
		if err != nil {
			var errDataNotFound *db.ErrDataNotFound[uint]

			if errors.As(err, &errDataNotFound) {
				return response.ErrNotFound(c, fmt.Errorf("employee ID %d could not be found", id))
			}

			return response.ErrInternalServer(c, fmt.Errorf("failed to find employee ID: %q", err.Error()))
		}

		return response.OK(c, employee)
	}
}

func handleDeleteEmployee(store store.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse employee ID: %v", err))
		}

		if err := store.DeleteEmployee(c.UserContext(), uint(id)); err != nil {
			return response.ErrInternalServer(c, fmt.Errorf("Failed to delete employee ID %d: %v", id, err))
		}

		return response.OK[any](c, nil)
	}
}

func handleUpdateEmployee(store store.Store, v *validator.Validator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("Unable to parse employee ID: %v", err))
		}

		var update models.UpsertEmployee
		if err := c.BodyParser(&update); err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("Unable to parse request body: %v", err))
		}

		err = v.Validate(update)
		if err != nil {
			return response.ErrUnprocessableEntity(c, v.SerializeErrors(err, update))
		}

		employee, err := store.UpdateEmployee(c.UserContext(), uint(id), update)
		if err != nil {
			return response.ErrInternalServer(c, fmt.Errorf("Failed to update employee ID %d: %v", id, err))
		}

		return response.OK(c, employee)
	}
}

func RegisterEmployee(router fiber.Router, store store.Store, v *validator.Validator) {
	router.Get("/", handleGetListEmployee(store))
	router.Post("/", handleCreateEmployee(store, v))
	router.Get("/:id", handleGetDetailEmployee(store))
	router.Delete("/:id", handleDeleteEmployee(store))
	router.Put("/:id", handleUpdateEmployee(store, v))
}
