package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/mattfan00/wfht/app"
	"github.com/mattfan00/wfht/store"
	_ "github.com/mattn/go-sqlite3"
)

const PORT = 8080
const DB_CONN = "/Users/matthewfan/sqlite/db/wfht.db"

func main() {
	db, err := sqlx.Connect("sqlite3", DB_CONN)
	if err != nil {
		panic(err)
	}
	log.Printf("connected to DB: %s\n", DB_CONN)

	eventStore := store.NewEventStore(db)

	a := app.New(eventStore)

	log.Printf("listening on port %d\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), a.Routes())
}
