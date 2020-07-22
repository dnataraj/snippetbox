package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, w http.ResponseWriter, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()

	// add flash messages, if any
	session, err := app.store.Get(r, "snippetbox-session")
	if err != nil {
		app.serverError(w, err)
		return nil
	}

	var flash = ""
	if f := session.Flashes(); len(f) > 0 {
		flash = f[0].(string)
	}
	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, err)
		return nil
	}
	td.Flash = flash

	td.IsAuthenticated = app.isAuthenticated(r)

	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, w, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (app *application) addFlash(w http.ResponseWriter, r *http.Request, msg string) error {
	session, err := app.store.Get(r, "snippetbox-session")
	if err != nil {
		return err
	}
	session.AddFlash(msg)
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	session, err := app.store.Get(r, "snippetbox-session")
	if err != nil {
		app.errorLog.Println("unable to read session", err)
		return false
	}
	_, ok := session.Values["authenticatedUserID"]
	return ok
}
