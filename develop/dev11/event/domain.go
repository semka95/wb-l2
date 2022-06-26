package event

import "time"

type User struct {
	ID     uint64  `json:"id,omitempty"`
	Events []Event `json:"events,omitempty"`
}

type Event struct {
	ID    uint64    `json:"id,omitempty"`
	Title string    `json:"title,omitempty"`
	Date  time.Time `json:"date,omitempty"`
}

type EventRepository interface {
	Create(user_id uint64, e Event) error
	Update(user_id uint64, e Event) error
	Delete(user_id uint64, event_id uint64) error
	GetForDay(user_id uint64, day time.Time) ([]Event, error)
	GetForWeek(user_id uint64, week time.Time) ([]Event, error)
	GetForMonth(user_id uint64, month time.Time) ([]Event, error)
}
