package holiday

import (
	"github.com/dwadp/attendance-api/models"
	"time"
)

const (
	Weekend int = iota
	NationalHoliday
)

type Service struct {
	holidays []*models.Holiday
}

func NewService(holidays []*models.Holiday) *Service {
	return &Service{holidays}
}

func (h *Service) IsHolidayExistOn(date time.Time) *models.Holiday {
	for _, item := range h.holidays {
		if item.Type == Weekend && item.Weekday != nil {
			if date.Weekday() == *item.Weekday {
				return item
			}
		} else if item.Type == NationalHoliday && item.Date != nil {
			hDate := item.Date.T
			newDate := time.Date(hDate.Year(), hDate.Month(), hDate.Day(), 0, 0, 0, 0, hDate.Location())

			if date.Equal(newDate) {
				return item
			}
		}
	}

	return nil
}
