package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const DB_CONN = "/Users/matthewfan/sqlite/db/wfht.db"

var db *sqlx.DB

func main() {
	if len(os.Args) != 2 {
		panic("incorrect num of args")
	}

	action := os.Args[1]

	db = sqlx.MustConnect("sqlite3", DB_CONN)
	log.Printf("connected to DB: %s\n", DB_CONN)

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
            date DATETIME NOT NULL,
            type INT NOT NULL,
            is_sys INT NOT NULL
        )
    `)
}

func mock() {
	db.MustExec("DELETE FROM event WHERE is_sys = true")

	ratio := 3.0 / 7.0
	numDaysSoFar := time.Now().YearDay()
	numDays := int(math.Ceil(ratio * float64(numDaysSoFar)))

	currYear := time.Now().Year()
	d := time.Date(currYear, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < numDays; i++ {
		db.MustExec(`
            INSERT INTO event (date, type, is_sys)
            VALUES ($1, $2, $3)
        `, d, 0, true)
		d = d.Add(time.Hour * 24)
	}

	fmt.Println(ratio)
	fmt.Println(numDaysSoFar)
	fmt.Println(numDays)
}
