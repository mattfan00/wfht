package store

import (
	"time"

	"github.com/jmoiron/sqlx"
)

var Events = []Event{
	{
		Date: time.Date(2023, 10, 15, 0, 0, 0, 0, time.Local),
		Type: "in",
	},
	{
		Date: time.Date(2023, 10, 16, 0, 0, 0, 0, time.Local),
		Type: "in",
	},
	{
		Date: time.Date(2023, 10, 17, 0, 0, 0, 0, time.Local),
		Type: "in",
	},
	{
		Date: time.Date(2023, 10, 18, 0, 0, 0, 0, time.Local),
		Type: "in",
	},
	{
		Date: time.Date(2023, 10, 19, 0, 0, 0, 0, time.Local),
		Type: "in",
	},
}

type Event struct {
	Date time.Time
	Type string
}

type EventStore struct {
	db *sqlx.DB
}

func NewEventStore(db *sqlx.DB) *EventStore {
	return &EventStore{
		db: db,
	}
}

func (e *EventStore) InsertMultiple(events []Event) {
	Events = append(Events, events...)
}
