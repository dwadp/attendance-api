package models

import (
	"github.com/dwadp/attendance-api/store/db"
	"time"
)

type Shift struct {
	ID        uint                `json:"id"`
	Name      string              `json:"name"`
	In        db.Time             `json:"in"`
	Out       db.Time             `json:"out"`
	IsDefault bool                `json:"is_default"`
	CreatedAt db.NullableDateTime `json:"created_at"`
	UpdatedAt db.NullableDateTime `json:"updated_at"`
}

func (s *Shift) GetIn() time.Time {
	now := time.Now()
	in := s.In.T

	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		in.Hour(),
		in.Minute(),
		in.Second(),
		0,
		now.Location(),
	)
}

func (s *Shift) GetOut() time.Time {
	now := time.Now()
	out := s.Out.T
	date := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		out.Hour(),
		out.Minute(),
		out.Second(),
		0,
		now.Location(),
	)

	if out.Before(s.In.T) {
		return date.Add(24 * time.Hour)
	}

	return date
}

type UpsertShift struct {
	Name      string  `json:"name" validate:"required,max=50"`
	In        db.Time `json:"in" validate:"required,time"`
	Out       db.Time `json:"out" validate:"required,time"`
	IsDefault bool    `json:"is_default"`
}
