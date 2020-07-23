package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	// common middleware
	stdMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynMiddleware := alice.New(csrf.Protect([]byte(app.secret)), app.authenticate)

	r := mux.NewRouter()
	r.Handle("/", dynMiddleware.ThenFunc(app.home)).Methods(http.MethodGet)
	r.Handle("/snippet/create", dynMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm)).Methods(http.MethodGet)
	r.Handle("/snippet/create", dynMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet)).Methods(http.MethodPost)
	r.Handle("/snippet/{id}", dynMiddleware.ThenFunc(app.showSnippet)).Methods(http.MethodGet)

	r.Handle("/user/signup", dynMiddleware.ThenFunc(app.signupUserForm)).Methods(http.MethodGet)
	r.Handle("/user/signup", dynMiddleware.ThenFunc(app.signupUser)).Methods(http.MethodPost)
	r.Handle("/user/login", dynMiddleware.ThenFunc(app.loginUserForm)).Methods(http.MethodGet)
	r.Handle("/user/login", dynMiddleware.ThenFunc(app.loginUser)).Methods(http.MethodPost)
	r.Handle("/user/logout", dynMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser)).Methods(http.MethodPost)

	r.HandleFunc("/ping", ping).Methods(http.MethodGet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer)).Methods(http.MethodGet)

	return stdMiddleware.Then(r)
}
