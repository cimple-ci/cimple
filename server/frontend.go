package server

import (
	"log"
	"net/http"

	"github.com/lukesmith/cimple/database"
	"github.com/lukesmith/cimple/web_application"
)

type agentPool interface {
	GetAgents() ([]*Agent, error)
}

type frontEnd struct {
	app *web_application.Application
}

func (fe *frontEnd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fe.app.Router.ServeHTTP(w, r)
}

func NewFrontend(db database.CimpleDatabase, agentPool *agentpool, buildQueue *buildQueue, addr string, logger *log.Logger) http.Handler {
	app := web_application.NewApplication(&web_application.ApplicationOptions{
		ViewsDirectory:  "./server/frontend/templates",
		AssetsDirectory: "./server/frontend/assets",
		Host:            addr,
	})

	app.Asset("/css/prism.css")
	app.Asset("/js/prism.js")
	app.Asset("/js/application.js")

	registerHome(app, db, agentPool)
	registerAgents(app, agentPool, logger)
	registerProjects(app, db)
	registerBuilds(app, db, buildQueue, logger)

	return &frontEnd{
		app: app,
	}
}
