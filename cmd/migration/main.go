package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mattfan00/wfht/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rickb777/date/v2"
)

var db *sqlx.DB

func main() {
	if len(os.Args) != 2 {
		panic("incorrect num of args")
	}

	configPath := flag.String("c", "./config.yaml", "path to config file")
	conf, err := config.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	action := os.Args[1]

	db = sqlx.MustConnect("sqlite3", conf.DbConn)
	log.Printf("connected to DB: %s\n", conf.DbConn)

	switch action {
	case "create":
		create()
	case "mock":
		mock()
	}
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
