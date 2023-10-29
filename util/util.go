package util

import "time"

const GOAL = 3

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	date, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}
	d.Time = date
	return nil
}

