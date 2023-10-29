package main

import (
	"log"
	"os"

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
	}
}

func create() {
	db.MustExec(`
        CREATE TABLE IF NOT EXISTS event (
            date DATETIME NOT NULL,
            type TEXT NOT NULL
        )
    `)
}
