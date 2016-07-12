package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lukesmith/cimple/database"
	"github.com/lukesmith/cimple/web_application"
	"time"
)

type buildsHandler struct {
	db database.CimpleDatabase
}

type buildModel struct {
	Id          string    `json:"id"`
	Date        time.Time `json:"date"`
	ProjectUrl  string    `json:"project_url"`
	BuildUrl    string    `json:"project_url"`
	BuildOutput string
}

func registerBuilds(app *web_application.Application, db database.CimpleDatabase) {
	handler := &buildsHandler{
		db: db,
	}

	app.Handle("/projects/{project_key}/builds/{key}", handler.getDetails).Methods("GET").Name("build")
}

func (h *buildsHandler) getDetails(app *web_application.Application, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)

	p, err := h.db.GetProject(params["project_key"])
	if err != nil {
		return nil, err
	}

	build, err := h.db.GetBuild(p.Name, params["key"])
	if err != nil {
		return nil, err
	}

	projectUrl, _ := app.Router.Get("project").URL("key", p.Name)
	buildUrl, _ := app.Router.Get("build").URL("project_key", p.Name, "key", build.Id)

	bo, err := build.GetOutput()
	if err != nil {
		return nil, err
	}

	return buildModel{
		Id:          build.Id,
		ProjectUrl:  projectUrl.Path,
		BuildUrl:    buildUrl.Path,
		BuildOutput: string(bo),
	}, nil
}
