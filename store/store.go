package store

import (
	"context"
	"github.com/dwadp/attendance-api/models"
	"time"
)

type Store interface {
	// Employee

	FindAllEmployees(ctx context.Context) ([]*models.Employee, error)
	CreateEmployee(ctx context.Context, employee models.UpsertEmployee) (*models.Employee, error)
	FindEmployeeByID(ctx context.Context, id uint) (*models.Employee, error)
	UpdateEmployee(ctx context.Context, id uint, employee models.UpsertEmployee) (*models.Employee, error)
	DeleteEmployee(ctx context.Context, id uint) error

	// Shift

	FindAllShifts(ctx context.Context) ([]*models.Shift, error)
	CreateShift(ctx context.Context, shift models.UpsertShift) (*models.Shift, error)
	FindShiftByID(ctx context.Context, id uint) (*models.Shift, error)
	UpdateShift(ctx context.Context, id uint, shift models.UpsertShift) (*models.Shift, error)
	DeleteShift(ctx context.Context, id uint) error
	FindEmployeeShift(ctx context.Context, employeeID uint, shiftID uint, date time.Time) (*models.EmployeeShift, error)
	SaveEmployeeShift(ctx context.Context, employeeShift models.EmployeeShift) (*models.EmployeeShift, error)

	// Holidays

	FindAllHolidays(ctx context.Context, date time.Time) ([]*models.Holiday, error)

	// Day Offs

	FindDayOff(ctx context.Context, employeeID uint, date time.Time) (*models.DayOff, error)
	SaveDayOff(ctx context.Context, dayOff models.DayOff) (*models.DayOff, error)
}
