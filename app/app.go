package app

import (
	"html/template"
	"path/filepath"

	"github.com/mattfan00/wfht/store"
)

type App struct {
	eventStore *store.EventStore
	templates  map[string]*template.Template
}

func New(eventStore *store.EventStore, templates map[string]*template.Template) *App {
	return &App{
		eventStore: eventStore,
		templates:  templates,
	}
}

func NewTemplates() (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}

	rootPath := "./ui/views"
	pages, err := filepath.Glob(filepath.Join(rootPath, "pages/*.html"))
	if err != nil {
		return map[string]*template.Template{}, err
	}

	for _, pagePath := range pages {
		name := filepath.Base(pagePath)
		t := template.New(name)

		t.ParseFiles(
			filepath.Join(rootPath, "base.html"),
			pagePath,
		)

		templates[name] = t
	}

	return templates, nil
}
