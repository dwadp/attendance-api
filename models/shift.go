package models

import (
	"github.com/dwadp/attendance-api/store/db"
)

type Shift struct {
	ID        uint                `json:"id"`
	Name      string              `json:"name"`
	In        db.Time             `json:"in"`
	Out       db.Time             `json:"out"`
	CreatedAt db.NullableDateTime `json:"created_at"`
	UpdatedAt db.NullableDateTime `json:"updated_at"`
}

type UpsertShift struct {
	Name string  `json:"name" validate:"required,max=50"`
	In   db.Time `json:"in" validate:"required,time"`
	Out  db.Time `json:"out" validate:"required,time"`
}
