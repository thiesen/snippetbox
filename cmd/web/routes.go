package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.home)))
	mux.Get("/snippets/new", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.createSnippetForm)))
	mux.Post("/snippets", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.createSnippet)))
	mux.Get("/snippets/:id", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.showSnippet)))

	fileServer := http.FileServer(http.Dir("./ui//static"))

	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
