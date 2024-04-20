package models

import (
	"github.com/dwadp/attendance-api/store/db"
	"time"
)

type Shift struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	In        time.Time      `json:"in"`
	Out       time.Time      `json:"out"`
	CreatedAt db.SqlNullTime `json:"created_at"`
	UpdatedAt db.SqlNullTime `json:"updated_at"`
}
