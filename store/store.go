package store

import (
	"context"
	"github.com/dwadp/attendance-api/models"
)

type Store interface {
	FindAllEmployees(ctx context.Context) ([]*models.Employee, error)
	CreateEmployee(ctx context.Context, employee models.UpsertEmployee) (*models.Employee, error)
	FindEmployeeByID(ctx context.Context, id uint) (*models.Employee, error)
	UpdateEmployee(ctx context.Context, id uint, employee models.UpsertEmployee) (*models.Employee, error)
	DeleteEmployee(ctx context.Context, id uint) error
}
