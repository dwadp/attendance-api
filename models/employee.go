package models

import (
	"github.com/dwadp/attendance-api/store/db"
)

type Employee struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Phone     string         `json:"phone"`
	CreatedAt db.SqlNullTime `json:"created_at"`
	UpdatedAt db.SqlNullTime `json:"updated_at"`
}

type UpsertEmployee struct {
	Name  string `json:"name" validate:"required,min=5,max=50"`
	Phone string `json:"phone" validate:"required,max=25"`
}
