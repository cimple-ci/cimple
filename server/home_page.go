package server

import (
	"net/http"

	"github.com/lukesmith/cimple/database"
	"github.com/lukesmith/cimple/web_application"
)

type homeHandler struct {
	db        database.CimpleDatabase
	agentPool agentPool
}

type homeModel struct {
	Hello    string
	Projects []*database.Project
	Agents   []*Agent
}

func registerHome(app *web_application.Application, db database.CimpleDatabase, agentPool agentPool) {
	handler := &homeHandler{
		db:        db,
		agentPool: agentPool,
	}

	app.Handle("/", handler.getIndex).Methods("GET")
}

func (h *homeHandler) getIndex(app *web_application.Application, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	agents, _ := h.agentPool.GetAgents()
	a := homeModel{
		Hello:    "World",
		Projects: h.db.GetProjects(),
		Agents:   agents,
	}

	return a, nil
}
