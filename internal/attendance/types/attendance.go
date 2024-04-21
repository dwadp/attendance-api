package types

import (
	"strconv"
)

type Type string

type Status string

const (
	ClockIn  Type = "clock_in"
	ClockOut      = "clock_out"
)

const (
	Alpha      Status = "alpha"
	Early             = "early"
	Late              = "late"
	NoClockIn         = "no_clock_in"
	NoClockOut        = "no_clock_out"
	Valid             = "valid"
)

func (s Status) MarshalJSON() ([]byte, error) {
	if s == "" {
		return []byte("null"), nil
	}

	quoted := strconv.Quote(string(s))

	return []byte(quoted), nil
}
