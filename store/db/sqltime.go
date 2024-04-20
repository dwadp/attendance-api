package db

import (
	"database/sql"
	"encoding/json"
	"time"
)

type SqlNullTime struct {
	sql.NullTime
}

func (s *SqlNullTime) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}

	result, err := json.Marshal(s.Time.Format(time.DateTime))

	if err != nil {
		return nil, err
	}

	return result, nil
}
