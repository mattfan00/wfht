package app

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/alexedwards/scs/v2"
	configPkg "github.com/mattfan00/wfht/config"
	"github.com/mattfan00/wfht/store"
)

type App struct {
	eventStore     *store.EventStore
	templates      map[string]*template.Template
	sessionManager *scs.SessionManager
	config         *configPkg.Config
}

func New(
	eventStore *store.EventStore,
	templates map[string]*template.Template,
	sessionManager *scs.SessionManager,
    config *configPkg.Config,
) *App {
	return &App{
		eventStore:     eventStore,
		templates:      templates,
		sessionManager: sessionManager,
        config: config,
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

func (a *App) renderTemplate(
	w http.ResponseWriter,
	template string,
	templateName string,
	data any,
) {
	t, ok := a.templates[template]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	err := t.ExecuteTemplate(buf, templateName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

func (a *App) renderPage(w http.ResponseWriter, template string, data any) {
	a.renderTemplate(w, template, "base", data)
}

func (a *App) renderErrorTemplate(w http.ResponseWriter, e error, status int) {
	a.renderError(w, "error", e, status)
}

func (a *App) renderErrorPage(w http.ResponseWriter, e error, status int) {
	a.renderError(w, "base", e, status)
}

func (a *App) renderError(w http.ResponseWriter, templateName string, e error, status int) {
	w.WriteHeader(status)
	a.renderTemplate(w, "error.html", templateName, map[string]any{
		"Error": e,
	})
}
