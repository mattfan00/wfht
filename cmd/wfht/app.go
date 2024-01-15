package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/jmoiron/sqlx"
	appPkg "github.com/mattfan00/wfht/app"
	"github.com/mattfan00/wfht/config"
	"github.com/mattfan00/wfht/store"
)

type appProgram struct {
	fs         *flag.FlagSet
	args       []string
	configPath string
}

func newAppProgram(args []string) *appProgram {
	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	configPath := fs.String("c", "./config.yaml", "path to config file")

	return &appProgram{
		fs:         fs,
		args:       args,
		configPath: *configPath,
	}
}

func (a *appProgram) parse() error {
	return a.fs.Parse(a.args)
}

func (a *appProgram) name() string {
	return a.fs.Name()
}

func (a *appProgram) run() error {
	conf, err := config.ReadFile(a.configPath)
	if err != nil {
		return err
	}

	db, err := sqlx.Connect("sqlite3", conf.DbConn)
	if err != nil {
		return err
	}
	log.Printf("connected to DB: %s\n", conf.DbConn)

	eventStore := store.NewEventStore(db)

	templates, err := appPkg.NewTemplates()
	if err != nil {
		return err
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Store = memstore.New()

	app := appPkg.New(
		eventStore,
		templates,
		sessionManager,
		conf,
	)

	log.Printf("listening on port %d\n", conf.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), app.Routes())

	return nil
}
