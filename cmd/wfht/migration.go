package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mattfan00/wfht/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rickb777/date/v2"
)

type migrationProgram struct {
	fs         *flag.FlagSet
	args       []string
	configPath string
}

func newMigrationProgram(args []string) *migrationProgram {
	fs := flag.NewFlagSet("migration", flag.ContinueOnError)
	configPath := fs.String("c", "./config.yaml", "path to config file")

	return &migrationProgram{
		fs:         fs,
		args:       args,
		configPath: *configPath,
	}
}

func (m *migrationProgram) parse() error {
	return m.fs.Parse(m.args)
}

func (m *migrationProgram) name() string {
	return m.fs.Name()
}

var db *sqlx.DB

func (m *migrationProgram) run() error {
	if len(m.args) != 1 {
		return fmt.Errorf("incorrect num of args")
	}

	conf, err := config.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	action := m.args[0]

	db = sqlx.MustConnect("sqlite3", conf.DbConn)
	fmt.Printf("connected to DB: %s\n", conf.DbConn)

	switch action {
	case "create":
		create()
	case "mock":
		mock()
	}

	return err
}

func create() {
	db.MustExec(`
        CREATE TABLE IF NOT EXISTS event (
            date DATE PRIMARY KEY,
            type INT NOT NULL,
            is_sys INT NOT NULL,
            updated_on DATETIME NOT NULL
        )
    `)
	fmt.Println("created 'event' table")
}

func mock() {
	db.MustExec("DELETE FROM event WHERE is_sys = true")

	currDay := date.Today()
	currTime := time.Now()

	ratio := 3.0 / 7.0
	numDaysSoFar := currDay.YearDay()
	numDays := int(math.Ceil(ratio * float64(numDaysSoFar)))

	currYear := currDay.Year()

	d := date.New(currYear, 1, 1)
	for i := 0; i < numDays; i++ {
		db.MustExec(`
            INSERT INTO event (date, type, is_sys, updated_on)
            VALUES ($1, $2, $3, $4)
        `, d, 0, true, currTime)
		d = d + 1
	}

	fmt.Println(ratio)
	fmt.Println(numDaysSoFar)
	fmt.Println(numDays)
}
