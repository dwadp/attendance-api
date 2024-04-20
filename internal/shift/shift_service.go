package shift

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dwadp/attendance-api/internal"
	"github.com/dwadp/attendance-api/internal/holiday"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/store"
	"github.com/dwadp/attendance-api/store/db"
)

type Service struct {
	store    store.Store
	hService *holiday.Service
}

func NewService(store store.Store, hService *holiday.Service) *Service {
	return &Service{store: store, hService: hService}
}

func (s *Service) AssignEmployee(ctx context.Context, assign models.AssignEmployeeShift) (*models.EmployeeShift, error) {
	existingShift, err := s.store.FindEmployeeShift(ctx, assign.EmployeeID, assign.Date.T)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if existingShift != nil {
		return existingShift, nil
	}

	existingDayOff, err := s.store.FindDayOff(ctx, assign.EmployeeID, assign.Date.T)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if existingDayOff != nil {
		return nil, internal.ErrDayOffExistsOnDate
	}

	employee, err := s.store.FindEmployeeByID(ctx, assign.EmployeeID)
	if err != nil {
		var errDataNotFound *db.ErrDataNotFound[uint]
		if errors.As(err, &errDataNotFound) {
			return nil, fmt.Errorf("employee ID %d could not be found", assign.EmployeeID)
		}
	}

	shift, err := s.store.FindShiftByID(ctx, assign.ShiftID)
	if err != nil {
		var errDataNotFound *db.ErrDataNotFound[uint]
		if errors.As(err, &errDataNotFound) {
			return nil, fmt.Errorf("shift ID %d could not be found", assign.ShiftID)
		}
	}

	if h := s.hService.IsHolidayExistOn(assign.Date.T); h != nil {
		if h.Type == holiday.Weekend {
			return nil, internal.ErrIsOnHoliday
		} else if h.Type == holiday.NationalHoliday {
			return nil, internal.ErrIsOnHoliday
		}
	}

	result, err := s.store.SaveEmployeeShift(ctx, models.EmployeeShift{
		EmployeeID: employee.ID,
		ShiftID:    shift.ID,
		Date:       assign.Date,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
