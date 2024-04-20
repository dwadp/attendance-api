package internal

import "errors"

var (
	ErrDayOffExistsOnDate = errors.New("there is a day off exists already on this date")
	ErrIsOnHoliday        = errors.New("could not proceed on holiday")
	ErrShiftExists        = errors.New("there is a shift exists already on this day")
)
