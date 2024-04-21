package handlers

import (
	"encoding/binary"
	"fmt"
	"github.com/dwadp/attendance-api/internal/attendance"
	"github.com/dwadp/attendance-api/internal/attendance/types"
	"github.com/dwadp/attendance-api/internal/holiday"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/server/response"
	"github.com/dwadp/attendance-api/server/validator"
	"github.com/dwadp/attendance-api/store"
	"github.com/gofiber/fiber/v2"
	"strconv"
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

func handleListAttendances(service *attendance.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		employeeID, err := c.ParamsInt("employee_id")
		if err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse employee_id: %v", err))
		}

		attendances, err := service.FindAllEmployeeAttendances(c.UserContext(), uint(employeeID))
		if err != nil {
			return response.ErrInternalServer(c, err)
		}

		return response.OK(c, attendances)
	}
}

func handleExportAttendanceList(service *attendance.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		employeeID, err := c.ParamsInt("employee_id")
		if err != nil {
			return response.ErrBadRequest(c, fmt.Errorf("unable to parse employee_id: %v", err))
		}

		employee, f, err := service.ExportAttendance(c.UserContext(), uint(employeeID))
		if err != nil {
			return response.ErrInternalServer(c, fmt.Errorf("unable to export attendance: %w", err))
		}

		buff, err := f.WriteToBuffer()
		if err != nil {
			return response.ErrInternalServer(c, err)
		}

		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=Attendance_%s.xlsx", employee.Name))
		c.Set("Content-Type", "application/octet-stream")
		c.Set("Content-Length", strconv.Itoa(binary.Size(buff)))

		return c.Send(buff.Bytes())
	}
}

func RegisterAttendance(router fiber.Router, store store.Store, v *validator.Validator) {
	holidayService := holiday.NewService(store)
	attendanceService := attendance.NewService(store, holidayService)

	router.Post("/clock-in", handleRequestClockIn(attendanceService, v))
	router.Post("/clock-out", handleRequestClockOut(attendanceService, v))
	router.Get("/:employee_id", handleListAttendances(attendanceService))
	router.Get("/:employee_id/export", handleExportAttendanceList(attendanceService))
}
