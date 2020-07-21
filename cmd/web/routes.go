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

	r.HandleFunc("/user/signup", app.signupUserForm).Methods(http.MethodGet)
	r.HandleFunc("/user/signup", app.signupUser).Methods(http.MethodPost)
	r.HandleFunc("/user/login", app.loginUserForm).Methods(http.MethodGet)
	r.HandleFunc("/user/login", app.loginUser).Methods(http.MethodPost)
	r.HandleFunc("/user/logout", app.logoutUser).Methods(http.MethodPost)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer)).Methods(http.MethodGet)

	return stdMiddleware.Then(r)
}
