package store

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rickb777/date/v2"
)

type EventType int

const (
	EventTypeCheckIn EventType = iota
	EventTypeOff               // disabled
	EventTypeNone
)

func (et EventType) IsValid() bool {
	switch et {
	case EventTypeCheckIn, EventTypeNone:
		return true
	}
	return false
}

var EventTypeMap = map[EventType]string{
	EventTypeCheckIn: "Check In",
	EventTypeNone:    "None",
}

var test time.Month

type Event struct {
	Date      date.Date `db:"date"`
	Type      EventType `db:"type"`
	IsSys     bool      `db:"is_sys"`
	UpdatedOn time.Time `db:"updated_on"`
	Display   bool
}

func (e *Event) IsCheckIn() bool {
	return e.Type == EventTypeCheckIn
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

func generateEventStmt(clause string) string {
	return fmt.Sprintf(`
        SELECT 
            date,
            type,
            is_sys
        FROM event
        %s
    `, clause)
}

func (es *EventStore) GetByCurrYear() ([]Event, error) {
	events := []Event{}

	stmt := generateEventStmt(`
        WHERE strftime('%Y', date) = $1
    `)

	args := []any{
		strconv.Itoa(time.Now().Year()),
	}
	err := es.db.Select(&events, stmt, args...)
	if err != nil {
		return []Event{}, err
	}

	return events, nil
}

func (es *EventStore) GetByYearMonth(year int, month time.Month) (map[date.Date]Event, error) {
	events := []Event{}

	stmt := generateEventStmt(`
        WHERE strftime('%Y', date) = $1 AND strftime('%m', date) = $2
    `)

	args := []any{
		strconv.Itoa(year),
		fmt.Sprintf("%02d", int(month)),
	}
	err := es.db.Select(&events, stmt, args...)
	if err != nil {
		return map[date.Date]Event{}, err
	}

	r := map[date.Date]Event{}
	for _, event := range events {
		r[event.Date] = event
	}

	return r, nil
}
