package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/thiesen/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)

		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)

		return
	}

	data := &templateData{Snippets: s}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)

		return
	}

	snippet, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)

		return
	} else if err != nil {
		app.serverError(w, err)

		return
	}

	data := &templateData{Snippet: snippet}

	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	app.infoLog.Printf("%v", snippet)

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", 405)

		return
	}

	title := "foo"
	content := "bar"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
