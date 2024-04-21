package internal

import "errors"

var (
	ErrDayOffExistsOnDate = errors.New("there is a day off exists already on this date")
	ErrIsOnHoliday        = errors.New("could not proceed on holiday")
	ErrShiftExists        = errors.New("there is a shift exists already on this day")

	ErrNationalHolidayShouldNotHaveWeekday = errors.New("type national holiday should not have a weekday")
	ErrWeekendShouldNotHaveDate            = errors.New("type weekday should not have a specific date")
)
