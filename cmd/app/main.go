package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/mattfan00/wfht/app"
	"github.com/mattfan00/wfht/config"
	"github.com/mattfan00/wfht/store"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configPath := flag.String("c", "./config.yaml", "path to config file")
	port := flag.Int("p", 8080, "port")
	flag.Parse()

	conf, err := config.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Connect("sqlite3", conf.DbConn)
	if err != nil {
		panic(err)
	}
	log.Printf("connected to DB: %s\n", conf.DbConn)

	eventStore := store.NewEventStore(db)

    templates, err := app.NewTemplates()
	if err != nil {
		panic(err)
	}

	a := app.New(eventStore, templates)

	log.Printf("listening on port %d\n", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), a.Routes())
}
