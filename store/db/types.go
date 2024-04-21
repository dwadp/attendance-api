package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type NullableDateTime struct {
	sql.NullTime
}

func (s *NullableDateTime) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}

	result, err := json.Marshal(s.Time.Format(time.DateTime))

	if err != nil {
		return nil, err
	}

	return result, nil
}

type Date struct {
	T     time.Time
	Valid bool
}

// Scan implements the [Scanner] interface.
func (d *Date) Scan(value any) error {
	if value == nil {
		*d = Date{T: time.Time{}, Valid: false}
		return nil
	}

	switch v := value.(type) {
	case string:
		r, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}

		*d = Date{T: r, Valid: true}
	case time.Time:
		*d = Date{T: time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location()), Valid: true}
	case []byte:
		r, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return nil
		}
		*d = Date{T: r, Valid: true}
	}

	return nil
}

// Value implements the [driver.Valuer] interface.
func (d Date) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}

	return d.T.Format("2006-01-02"), nil
}

// MarshalJSON implements the [json.Marshaler] interface.
func (d *Date) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}

	result, err := json.Marshal(d.T.Format("2006-01-02"))

	if err != nil {
		return nil, err
	}

	return result, nil
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (d *Date) UnmarshalJSON(data []byte) error {
	if d == nil {
		return errors.New("UnmarshalJSON on nil pointer")
	}

	value := strings.Trim(string(data), "\"")

	res, err := time.Parse("2006-01-02", value)
	if err != nil {
		return fmt.Errorf("Unable to unmarshall time invalid format: %w", err)
	}

	*d = Date{T: res, Valid: true}
	return nil
}

func (d *Date) String() string {
	if !d.Valid {
		return ""
	}

	return d.T.Format("15:04:05")
}

// Time represents a 24 hours format time (HH:MM:SS) without a date. It could also be a null time
type Time struct {
	T     time.Time
	Valid bool
}

// Scan implements the [Scanner] interface.
func (t *Time) Scan(value any) error {
	if value == nil {
		*t = Time{T: time.Time{}, Valid: false}
		return nil
	}

	switch v := value.(type) {
	case string:
		r, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		*t = Time{T: r, Valid: true}
	case time.Time:
		*t = Time{T: v, Valid: true}
	case []byte:
		r, err := time.Parse("15:04:05", string(v))
		if err != nil {
			return nil
		}
		*t = Time{T: r, Valid: true}
	}

	return nil
}

// Value implements the [driver.Valuer] interface.
func (t Time) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}

	return t.T.Format("15:04:05"), nil
}

// MarshalJSON implements the [json.Marshaler] interface.
func (t *Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}

	result, err := json.Marshal(t.T.Format("15:04:05"))

	if err != nil {
		return nil, err
	}

	return result, nil
}

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (t *Time) UnmarshalJSON(data []byte) error {
	if t == nil {
		return errors.New("UnmarshalJSON on nil pointer")
	}

	value := strings.Trim(string(data), "\"")

	if len(strings.Split(value, ":")) == 2 {
		value = strings.TrimRight(value, ":")
		value += ":" + "00"
	}

	res, err := time.Parse("15:04:05", value)
	if err != nil {
		return fmt.Errorf("Unable to unmarshall time invalid format: %w", err)
	}

	*t = Time{T: res, Valid: true}
	return nil
}

func (t *Time) String() string {
	if !t.Valid {
		return ""
	}

	return t.T.Format("15:04:05")
}

type NullableInt64 struct {
	sql.NullInt64
}

func NewNullableInt64(i int64) NullableInt64 {
	return NullableInt64{sql.NullInt64{Int64: i, Valid: true}}
}

func (s *NullableInt64) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}

	result, err := json.Marshal(s.Int64)

	if err != nil {
		return nil, err
	}

	return result, nil
}

type NullableString struct {
	sql.NullString
}

func NewNullableString(s string) NullableString {
	return NullableString{sql.NullString{String: s, Valid: true}}
}

func (s *NullableString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}

	result, err := json.Marshal(s.String)

	if err != nil {
		return nil, err
	}

	return result, nil
}
