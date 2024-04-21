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
	"github.com/golang-module/carbon/v2"
	"github.com/xuri/excelize/v2"
	"slices"
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

func (s *Service) FindAllEmployeeAttendances(ctx context.Context, employeeID uint) ([]*models.Attendance, error) {
	attendances, err := s.store.FindAllAttendances(ctx, employeeID)
	if err != nil {
		return nil, err
	}

	now := carbon.Now().StartOfMonth()
	result := make([]*models.Attendance, now.DaysInMonth())

	for day := 1; day <= now.DaysInMonth(); day++ {
		date := now.SetDay(day)
		key := day - 1

		index := slices.IndexFunc(attendances, func(a *models.Attendance) bool {
			return date.ToDateString() == a.Date.T.Format("2006-01-02")
		})

		if index >= 0 {
			result[key] = attendances[index]
		} else {
			result[key] = &models.Attendance{
				ID:             0,
				EmployeeID:     employeeID,
				ClockInStatus:  types.Alpha,
				ClockOutStatus: types.Alpha,
				Date: db.Date{
					T:     date.StdTime(),
					Valid: true,
				},
				CreatedAt: date.StdTime(),
				UpdatedAt: date.StdTime(),
			}
		}
	}

	return result, nil
}

func (s *Service) ExportAttendance(ctx context.Context, employeeID uint) (*models.Employee, *excelize.File, error) {
	employee, err := s.store.FindEmployeeByID(ctx, employeeID)
	if err != nil {
		return nil, nil, err
	}

	attendances, err := s.FindAllEmployeeAttendances(ctx, employeeID)
	if err != nil {
		return nil, nil, err
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"
	sheet, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, nil, err
	}

	f.SetCellValue(sheetName, "A2", "Date")
	f.SetCellValue(sheetName, "B2", "Shift Name")
	f.SetCellValue(sheetName, "C2", "Shift In")
	f.SetCellValue(sheetName, "D2", "Shift Out")
	f.SetCellValue(sheetName, "E2", "Clock In")
	f.SetCellValue(sheetName, "F2", "Clock Out")
	f.SetCellValue(sheetName, "G2", "Clock In Status")
	f.SetCellValue(sheetName, "H2", "Clock Out Status")

	row := 3
	for _, attendance := range attendances {
		shiftName := "-"
		shiftIn := "-"
		shiftOut := "-"
		clockIn := "-"
		clockOut := "-"

		if attendance.ShiftName.Valid {
			shiftName = attendance.ShiftName.String
		}

		if !attendance.ShiftIn.IsZero() {
			shiftIn = attendance.ShiftIn.Format("2006-01-02 15:04-05")
		}

		if !attendance.ShiftOut.IsZero() {
			shiftOut = attendance.ShiftOut.Format("2006-01-02 15:04-05")
		}

		if attendance.ClockIn.Valid {
			clockIn = attendance.ClockIn.Time.Format("2006-01-02 15:04:05")
		}

		if attendance.ClockOut.Valid {
			clockOut = attendance.ClockOut.Time.Format("2006-01-02 15:04:05")
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), attendance.Date.T.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), shiftName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), shiftIn)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), shiftOut)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), clockIn)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), clockOut)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), attendance.GetClockInStatus())
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), attendance.GetClockOutStatus())

		row++
	}

	f.SetActiveSheet(sheet)

	return employee, f, nil
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
