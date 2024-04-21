package store

import (
	"context"
	holidayTypes "github.com/dwadp/attendance-api/internal/holiday/types"
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
	FindDefaultShift(ctx context.Context) (*models.Shift, error)
	FindEmployeeShift(ctx context.Context, employeeID uint, date time.Time) (*models.EmployeeShift, error)
	SaveEmployeeShift(ctx context.Context, employeeShift models.EmployeeShift) (*models.EmployeeShift, error)
	DeleteEmployeeShift(ctx context.Context, unassign models.UnassignEmployeeShift) error

	// Holidays

	FindHolidaysInDate(ctx context.Context, date time.Time) ([]*models.Holiday, error)
	FindAllHoliday(ctx context.Context) ([]*models.Holiday, error)
	FindHolidayByDateOrWeekday(ctx context.Context, date time.Time, weekday holidayTypes.Weekday) (*models.Holiday, error)
	FindHolidayByID(ctx context.Context, id uint) (*models.Holiday, error)
	CreateHoliday(ctx context.Context, holiday models.UpsertHoliday) (*models.Holiday, error)
	UpdateHoliday(ctx context.Context, id uint, holiday models.UpsertHoliday) (*models.Holiday, error)
	DeleteHoliday(ctx context.Context, id uint) error

	// Day Offs

	FindDayOff(ctx context.Context, employeeID uint, date time.Time) (*models.DayOff, error)
	SaveDayOff(ctx context.Context, dayOff models.DayOff) (*models.DayOff, error)

	// Attendances

	SaveAttendance(ctx context.Context, attendance models.Attendance) (*models.Attendance, error)
	UpdateAttendance(ctx context.Context, attendance *models.Attendance) (*models.Attendance, error)
	FindAttendanceByEmployeeID(ctx context.Context, employeeID uint, date time.Time) (*models.Attendance, error)
	FindAllAttendances(ctx context.Context, employeeID uint) ([]*models.Attendance, error)
}
