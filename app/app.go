package app

import (
	"github.com/mattfan00/wfht/store"
)

type App struct {
	eventStore *store.EventStore
}

func New(eventStore *store.EventStore) *App {
	return &App{
		eventStore: eventStore,
	}
}
