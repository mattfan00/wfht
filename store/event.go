package store

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type EventType int

const (
	EventTypeCheckIn EventType = iota
	EventTypeOff
)

func (et EventType) IsValid() bool {
	switch et {
	case EventTypeCheckIn, EventTypeOff:
		return true
	}
	return false
}

type Event struct {
	Date      time.Time `db:"date"`
	Type      EventType `db:"type"`
	IsSys     bool      `db:"is_sys"`
	UpdatedOn time.Time `db:"updated_on"`
}

type EventStore struct {
	db *sqlx.DB
}

func NewEventStore(db *sqlx.DB) *EventStore {
	return &EventStore{
		db: db,
	}
}

func (es *EventStore) UpsertMultiple(events []Event) error {
	tx, err := es.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	currTime := time.Now()
	for _, event := range events {
		log.Printf("upserting event: %+v", event)
		stmt := `
            INSERT INTO event (date, type, is_sys, updated_on)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (date) DO UPDATE SET
                type = excluded.type,
                is_sys = excluded.is_sys,
                updated_on = excluded.updated_on
        `
		args := []any{
			event.Date,
			event.Type,
			event.IsSys,
			currTime,
		}
		_, err = tx.Exec(stmt, args...)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	return err
}

func (es *EventStore) GetByCurrYear() ([]Event, error) {
	events := []Event{}

	stmt := `
        SELECT 
            date,
            type,
            is_sys
        FROM event
        WHERE strftime('%Y', $1)
    `
	args := []any{time.Now().Year()}
	err := es.db.Select(&events, stmt, args...)
	if err != nil {
		return []Event{}, err
	}

	return events, nil
}
