package app

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

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

type HomeData struct {
	SubmitEventType store.EventType
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

	events, err := a.eventStore.GetByCurrYear()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	numCheckIn := 0
	numOff := 0
	checkedInToday := false
	currTime := time.Now()
	currDate := time.Date(currTime.Year(), currTime.Month(), currTime.Day(), 0, 0, 0, 0, time.UTC)
	for _, event := range events {
		if event.Type == store.EventTypeCheckIn {
			numCheckIn++
		} else if event.Type == store.EventTypeOff {
			numOff++
		}

		if event.Date.Equal(currDate) {
			checkedInToday = true
		}
	}

	t.Execute(w, map[string]any{
		"EventTypeCheckIn": store.EventTypeCheckIn,
		"CheckedInToday":   checkedInToday,
		"NumCheckIn":       numCheckIn,
		"NumOff":           numOff,
		"CurrDate":         currDate.Format("2006-01-02"),
	})
}

type EventRequest struct {
	Dates []util.Date     `json:"dates"`
	Type  store.EventType `json:"type"`
}

func (a *App) createEvents(w http.ResponseWriter, r *http.Request) {
	var req EventRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !req.Type.IsValid() {
		http.Error(w, "invalid event type", http.StatusInternalServerError)
		return
	}

	newEvents := []store.Event{}
	for _, date := range req.Dates {
		newEvents = append(newEvents, store.Event{
			Date:  date.Time,
			Type:  req.Type,
			IsSys: false,
		})
	}
	a.eventStore.UpsertMultiple(newEvents)

	w.Header().Add("HX-Redirect", "/")
	w.Write(nil)
}
