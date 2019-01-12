package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/thiesen/snippetbox/pkg/forms"
	"github.com/thiesen/snippetbox/pkg/models"
)

type templateData struct {
	AuthenticatedUser int
	CSRFToken         string
	CurrentYear       int
	Flash             string
	Form              *forms.Form
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/*.page.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/*.layout.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/*.partial.tmpl")
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
