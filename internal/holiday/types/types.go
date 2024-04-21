package types

import "strconv"

type Holiday int

const (
	Weekend Holiday = iota
	NationalHoliday
)

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	None
)

func (w *Weekday) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || len(data) == 0 {
		*w = None
		return nil
	}

	d, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}

	switch d {
	case 0:
		*w = Sunday
		break
	case 1:
		*w = Monday
		break
	case 2:
		*w = Tuesday
		break
	case 3:
		*w = Wednesday
		break
	case 4:
		*w = Thursday
		break
	case 5:
		*w = Friday
		break
	case 6:
		*w = Saturday
	case 7:
		*w = None
	}

	return nil
}
