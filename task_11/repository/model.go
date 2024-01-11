package repository

import (
	"fmt"
	"strings"
	"time"
)

type Event struct {
	EventId     int        `json:"event_id"`
	UserId      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Date        CustomDate `json:"date"`
}

type CustomDate struct {
	time.Time
}

func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return
	}
	parsedTime, err := time.Parse("2006-01-02T15:04", s)
	if err != nil {
		parsedTime, err = time.Parse("2006-01-02T15:04:00Z", s)
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02", s)
			if err != nil {
				return fmt.Errorf("date format: e.g. 2022-05-10T14:10 error: %v", err)
			}
		}
	}
	*c = CustomDate{parsedTime}
	return nil
}
