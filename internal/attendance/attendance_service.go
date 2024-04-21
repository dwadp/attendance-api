package attendance

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dwadp/attendance-api/internal"
	"github.com/dwadp/attendance-api/internal/attendance/types"
	"github.com/dwadp/attendance-api/internal/holiday"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/store"
	"github.com/dwadp/attendance-api/store/db"
	"time"
)

type Service struct {
	store          store.Store
	holidayService *holiday.Service
}

func NewService(store store.Store, holidayService *holiday.Service) *Service {
	return &Service{store, holidayService}
}

func (s *Service) ClockIn(ctx context.Context, req models.AttendanceRequest) (*models.Attendance, error) {
	now := time.Now()

	existing, err := s.store.FindAttendanceByEmployeeID(ctx, req.EmployeeID, now)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if existing != nil {
		return existing, nil
	}

	attendance, err := s.createNewAttendance(ctx, req, now, types.ClockIn)
	if err != nil {
		return nil, err
	}

	return s.store.SaveAttendance(ctx, *attendance)
}

func (s *Service) ClockOut(ctx context.Context, req models.AttendanceRequest) (*models.Attendance, error) {
	now := time.Now()
	existing, err := s.store.FindAttendanceByEmployeeID(ctx, req.EmployeeID, now)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if existing != nil {
		if !existing.ClockOut.Valid && existing.ClockOutStatus == "" {
			existing.ClockOut = db.NullableDateTime{NullTime: sql.NullTime{
				Time:  now,
				Valid: true,
			}}
			existing.ClockOutStatus = s.status(now, existing.ShiftOut)

			return s.store.UpdateAttendance(ctx, existing)
		}

		return existing, nil
	}

	attendance, err := s.createNewAttendance(ctx, req, now, types.ClockOut)
	if err != nil {
		return nil, err
	}

	return s.store.SaveAttendance(ctx, *attendance)
}

func (s *Service) getShift(ctx context.Context, employeeID uint, date time.Time) (shift *models.Shift, err error) {
	employeeShift, err := s.store.FindEmployeeShift(ctx, employeeID, date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shift, err = s.store.FindDefaultShift(ctx)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if employeeShift != nil {
		shift, err = s.store.FindShiftByID(ctx, employeeShift.ShiftID)
	}

	if shift == nil {
		return nil, fmt.Errorf("unable to request attendance because there is no shift")
	}

	return shift, nil
}

func (s *Service) createNewAttendance(ctx context.Context, req models.AttendanceRequest, date time.Time, t types.Type) (*models.Attendance, error) {
	// Check if the employee data is valid
	_, err := s.store.FindEmployeeByID(ctx, req.EmployeeID)
	if err != nil {
		var notFound *db.ErrDataNotFound[uint]
		if errors.As(err, &notFound) {
			return nil, notFound
		}
		return nil, err
	}

	dayOff, err := s.store.FindDayOff(ctx, req.EmployeeID, date)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if dayOff != nil {
		return nil, internal.ErrDayOffExistsOnDate
	}

	existingHoliday, err := s.holidayService.IsHolidayExistOn(ctx, date)
	if err != nil {
		return nil, err
	}

	if existingHoliday != nil {
		return nil, internal.ErrIsOnHoliday
	}

	shift, err := s.getShift(ctx, req.EmployeeID, date)
	if err != nil {
		return nil, err
	}

	attendance := models.Attendance{
		EmployeeID: req.EmployeeID,
		ShiftID:    db.NewNullableInt64(int64(shift.ID)),
		ShiftName:  db.NewNullableString(shift.Name),
		ShiftIn:    shift.GetIn(),
		ShiftOut:   shift.GetOut(),
		Date: db.Date{
			T:     date,
			Valid: true,
		},
	}

	switch t {
	case types.ClockIn:
		attendance.ClockIn = db.NullableDateTime{NullTime: sql.NullTime{
			Time:  date,
			Valid: true,
		}}
		attendance.ClockInStatus = s.status(date, attendance.ShiftIn)
	case types.ClockOut:
		attendance.ClockOut = db.NullableDateTime{NullTime: sql.NullTime{
			Time:  date,
			Valid: true,
		}}
		attendance.ClockOutStatus = s.status(date, attendance.ShiftOut)
	}

	return &attendance, nil
}

func (s *Service) status(now time.Time, tm time.Time) types.Status {
	now = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		0,
		0,
		now.Location(),
	)

	if now.Before(tm) {
		return types.Early
	} else if now.After(tm) {
		return types.Late
	}

	return types.Valid
}
