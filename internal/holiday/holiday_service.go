package holiday

import (
	"context"
	"github.com/dwadp/attendance-api/models"
	"github.com/dwadp/attendance-api/store"
	"time"
)

const (
	Weekend int = iota
	NationalHoliday
)

type Service struct {
	store store.Store
}

func NewService(store store.Store) *Service {
	return &Service{store}
}

func (h *Service) IsHolidayExistOn(ctx context.Context, date time.Time) (*models.Holiday, error) {
	holidays, err := h.store.FindAllHolidays(ctx, date)
	if err != nil {
		return nil, err
	}

	for _, item := range holidays {
		if item.Type == Weekend && item.Weekday != nil {
			if date.Weekday() == *item.Weekday {
				return item, nil
			}
		} else if item.Type == NationalHoliday && item.Date != nil {
			hDate := item.Date.T
			newDate := time.Date(hDate.Year(), hDate.Month(), hDate.Day(), 0, 0, 0, 0, hDate.Location())

			if date.Equal(newDate) {
				return item, nil
			}
		}
	}

	return nil, nil
}
