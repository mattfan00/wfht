package app

import (
	"encoding/json"
	"html/template"
	"log"
	"math"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mattfan00/wfht/store"
	"github.com/rickb777/date/v2"
)

func (a *App) Routes() *chi.Mux {
	router := chi.NewRouter()

	staticFileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handle("/static/*", http.StripPrefix("/static/", staticFileServer))

	router.Get("/", a.getHomePage)
	router.Get("/calendar", a.getCalendarPage)
	router.Post("/events", a.createEvents)

	return router
}

func (a *App) getHomePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./ui/views/base.html",
		"./ui/views/pages/home.html",
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
	checkedInToday := false
	currDate := date.Today()
	for _, event := range events {
		if event.Type == store.EventTypeCheckIn {
			numCheckIn++
		}

		if event.Date == currDate {
			checkedInToday = true
		}
	}

	numDaysSoFar := currDate.YearDay()
	currRatio := float64(numCheckIn) / float64(numDaysSoFar)
	currAvgCheckIn := currRatio * 7
	numDaysGoal := math.Ceil(365 * (3.0 / 7.0))

	t.Execute(w, map[string]any{
		"EventTypeCheckIn": store.EventTypeCheckIn,
		"CheckedInToday":   checkedInToday,
		"CurrAvgCheckIn":   currAvgCheckIn,
		"NumDaysGoal":      numDaysGoal,
		"NumCheckIn":       numCheckIn,
	})
}

func sameMonth(d1 date.Date, d2 date.Date) bool {
	return d1.Month() == d2.Month()
}

func (a *App) getCalendarPage(w http.ResponseWriter, r *http.Request) {
	t := template.New("")

	t.Funcs(template.FuncMap{
		"sameMonth": sameMonth,
	})

	t, err := t.ParseFiles(
		"./ui/views/base.html",
		"./ui/views/pages/calendar.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currDate := date.Today()

	eventMap, err := a.eventStore.GetByYearMonth(currDate.Year(), currDate.Month())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	firstOfMonthDate := date.New(currDate.Year(), currDate.Month(), 1)
	lastOfMonthDate := firstOfMonthDate.AddDate(0, 1, -1)

	calendar := []store.Event{}

	// previous month
	for i := int(firstOfMonthDate.Weekday()); i > 0; i-- {
		calendar = append(calendar, store.Event{
			Date: firstOfMonthDate.AddDate(0, 0, -i),
			Type: store.EventTypeNone,
		})
	}

	// current month
	for i := 0; i < lastOfMonthDate.Day(); i++ {
		var e store.Event
		d := firstOfMonthDate.AddDate(0, 0, i)
		val, ok := eventMap[d]
		if ok {
			e = val
		} else {
			e = store.Event{
				Date: d,
				Type: store.EventTypeNone,
			}
		}
		calendar = append(calendar, e)
	}

	// next month
	for i := 1; i < 7-int(lastOfMonthDate.Weekday()); i++ {
		calendar = append(calendar, store.Event{
			Date: lastOfMonthDate.AddDate(0, 0, i),
			Type: store.EventTypeNone,
		})
	}

	t.ExecuteTemplate(w, "base.html", map[string]any{
		"Calendar":         calendar,
		"CurrDate":         currDate,
		"EventTypeCheckIn": store.EventTypeCheckIn,
	})
}

type EventRequest struct {
	Dates []date.Date     `json:"dates"`
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

    log.Printf("%+v", req)

	newEvents := []store.Event{}
	for _, date := range req.Dates {
		newEvents = append(newEvents, store.Event{
			Date:  date,
			Type:  req.Type,
			IsSys: false,
		})
	}

	err = a.eventStore.UpsertMultiple(newEvents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("HX-Redirect", "/")
	w.Write(nil)
}
