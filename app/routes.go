package app

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mattfan00/wfht/store"
	"github.com/mattfan00/wfht/util"
)

func (a *App) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", a.getHome)
	router.Post("/events", a.createEvents)

	return router
}

func (a *App) getHome(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./public/views/base.html",
		"./public/views/pages/home.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Print("hit home")

	t.Execute(w, nil)
}

type EventRequest struct {
	Dates []util.Date `json:"dates"`
	Type  string      `json:"type"`
}

func (a *App) createEvents(w http.ResponseWriter, r *http.Request) {
	var req EventRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newEvents := []store.Event{}
	for _, date := range req.Dates {
		newEvents = append(newEvents, store.Event{
			Date: date.Time,
			Type: req.Type,
		})

	}
	a.eventStore.InsertMultiple(newEvents)

	http.Redirect(w, r, "/", http.StatusOK)
}
