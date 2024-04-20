package dayoff

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dwadp/attendance-api/internal"
	"github.com/dwadp/attendance-api/internal/holiday"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/store"
)

type Service struct {
	store          store.Store
	holidayService *holiday.Service
}

func NewService(store store.Store, holidayService *holiday.Service) *Service {
	return &Service{
		store:          store,
		holidayService: holidayService,
	}
}

func (s *Service) Create(ctx context.Context, request models.DayOffRequest) (*models.DayOff, error) {
	if h := s.holidayService.IsHolidayExistOn(request.Date.T); h != nil {
		return nil, internal.ErrIsOnHoliday
	}

	existingShift, err := s.store.FindEmployeeShift(ctx, request.EmployeeID, request.Date.T)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if existingShift != nil {
		return nil, internal.ErrShiftExists
	}

	existingDayOff, err := s.store.FindDayOff(ctx, request.EmployeeID, request.Date.T)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if existingDayOff != nil {
		return existingDayOff, nil
	}

	return s.store.SaveDayOff(ctx, models.DayOff{
		EmployeeID:  request.EmployeeID,
		Description: request.Description,
		Date:        request.Date,
	})
}
