package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mattfan00/wfht/util"
)

type App struct {
}

func New() *App {
	return &App{}
}

func (a *App) Routes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		t.Execute(w, "hi")
	})

	router.Post("/api/event", a.createEvent)

	return router
}

type EventRequest struct {
	Date util.Date `json:"date"`
	Type string    `json:"type"`
}

func (a *App) createEvent(w http.ResponseWriter, r *http.Request) {
	var req EventRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	fmt.Fprintf(w, "req: %+v", req)
}
