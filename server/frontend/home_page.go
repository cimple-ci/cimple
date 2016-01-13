package frontend

import (
	"net/http"

	"github.com/lukesmith/cimple/database"
	"github.com/lukesmith/cimple/web_application"
)

type homeHandler struct {
	db database.CimpleDatabase
}

type homeModel struct {
	Hello    string
	Projects []*database.Project
}

func registerHome(app *web_application.Application, db database.CimpleDatabase) {
	handler := &homeHandler{
		db: db,
	}

	app.Handle("/", handler.getIndex).Methods("GET")
}

func (h *homeHandler) getIndex(app *web_application.Application, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	a := homeModel{
		Hello:    "World",
		Projects: h.db.GetProjects(),
	}

	return a, nil
}
