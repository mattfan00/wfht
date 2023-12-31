package app

import (
	"encoding/json"
	"html/template"
	"math"
	"net/http"
	"time"

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
	router.Get("/calendar/partial", a.getCalendarPartial)
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

			if event.Date == currDate {
				checkedInToday = true
			}
		}

	}

	numDaysSoFar := currDate.YearDay()
	currRatio := float64(numCheckIn) / float64(numDaysSoFar)
	currAvgCheckIn := currRatio * 7
	numDaysGoal := math.Ceil(365 * (3.0 / 7.0))

	t.ExecuteTemplate(w, "base", map[string]any{
		"EventTypeCheckIn": store.EventTypeCheckIn,
		"CheckedInToday":   checkedInToday,
		"CurrAvgCheckIn":   currAvgCheckIn,
		"NumDaysGoal":      numDaysGoal,
		"NumCheckIn":       numCheckIn,
	})
}

type CalendarOption struct {
	Month      time.Month
	Value      date.Date
	IsSelected bool
}

func (a *App) getCalendarPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./ui/views/base.html",
		"./ui/views/pages/calendar.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currDate := date.Today()

	data, err := a.generateCalendarPartialData(currDate.Year(), currDate.Month())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	calendarOptions := []CalendarOption{}
	for i := time.January; i <= time.December; i++ {
		calendarOptions = append(calendarOptions, CalendarOption{
			Month:      i,
			Value:      date.New(currDate.Year(), i, 1),
			IsSelected: i == currDate.Month(),
		})
	}

	data["CalendarOptions"] = calendarOptions

	t.ExecuteTemplate(w, "base", data)
}

func (a *App) getCalendarPartial(w http.ResponseWriter, r *http.Request) {
	monthParam := r.URL.Query().Get("month")
	d, err := date.ParseISO(monthParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles(
		"./ui/views/pages/calendar.html",
	)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := a.generateCalendarPartialData(d.Year(), d.Month())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.ExecuteTemplate(w, "calendar", data)
}

func (a *App) generateCalendarPartialData(year int, month time.Month) (map[string]any, error) {
	eventMap, err := a.eventStore.GetByYearMonth(year, month)
	if err != nil {
		return map[string]any{}, err
	}

	firstOfMonthDate := date.New(year, month, 1)
	lastOfMonthDate := firstOfMonthDate.AddDate(0, 1, -1)

	calendar := []store.Event{}

	// previous month
	for i := int(firstOfMonthDate.Weekday()); i > 0; i-- {
		calendar = append(calendar, store.Event{
			Date:    firstOfMonthDate.AddDate(0, 0, -i),
			Type:    store.EventTypeNone,
			Display: false,
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
		e.Display = true
		calendar = append(calendar, e)
	}

	// next month
	for i := 1; i < 7-int(lastOfMonthDate.Weekday()); i++ {
		calendar = append(calendar, store.Event{
			Date:    lastOfMonthDate.AddDate(0, 0, i),
			Type:    store.EventTypeNone,
			Display: false,
		})
	}

	data := map[string]any{
		"Calendar":         calendar,
		"EventTypeCheckIn": store.EventTypeCheckIn,
		"EventTypeNone":    store.EventTypeNone,
	}

	return data, nil
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
