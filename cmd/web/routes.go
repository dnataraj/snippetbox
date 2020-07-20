package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	// common middleware
	stdMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	r := mux.NewRouter()
	r.HandleFunc("/", app.home).Methods(http.MethodGet)
	r.HandleFunc("/snippet/create", app.createSnippetForm).Methods(http.MethodGet)
	r.HandleFunc("/snippet/create", app.createSnippet).Methods(http.MethodPost)
	r.HandleFunc("/snippet/{id}", app.showSnippet).Methods(http.MethodGet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server
	r.Handle("/static/", http.StripPrefix("/static", fileServer)).Methods(http.MethodGet)

	return stdMiddleware.Then(r)
}
