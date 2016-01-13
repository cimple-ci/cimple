package frontend

import (
	"net/http"

	"github.com/lukesmith/cimple/database"
	"github.com/lukesmith/cimple/web_application"
)

type frontEnd struct {
	app *web_application.Application
}

func (fe *frontEnd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fe.app.Router.ServeHTTP(w, r)
}

func NewFrontend(db database.CimpleDatabase) http.Handler {
	app := web_application.NewApplication(&web_application.ApplicationOptions{
		ViewsDirectory: "frontend/templates",
	})

	app.Static("/assets/js/application.js")

	registerHome(app, db)
	registerProjects(app, db)

	return &frontEnd{
		app: app,
	}
}
