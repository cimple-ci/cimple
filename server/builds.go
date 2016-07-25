package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lukesmith/cimple/database"
	"github.com/lukesmith/cimple/messages"
	"github.com/lukesmith/cimple/web_application"
	"log"
	"time"
)

type buildsHandler struct {
	db         database.CimpleDatabase
	buildQueue BuildQueue
	logger     *log.Logger
}

type buildModel struct {
	Id          string    `json:"id"`
	Date        time.Time `json:"date"`
	ProjectUrl  string    `json:"project_url"`
	BuildUrl    string    `json:"project_url"`
	BuildOutput string
}

type submitBuildModel struct {
	Url    string
	Commit string
}

func registerBuilds(app *web_application.Application, db database.CimpleDatabase, buildQueue BuildQueue, logger *log.Logger) {
	handler := &buildsHandler{
		db:         db,
		buildQueue: buildQueue,
		logger:     logger,
	}

	app.Handle("/projects/{project_key}/builds/{key}", handler.getDetails).Methods("GET").Name("build")
	app.Handle("/builds", handler.submitBuild).Methods("POST").Name("submitBuild")
}

func (h *buildsHandler) submitBuild(app *web_application.Application, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	var submitModel submitBuildModel
	err := decoder.Decode(&submitModel)
	if err != nil {
		w.Write([]byte("Unprocessible entity"))
		return nil, err
	} else {
		w.WriteHeader(http.StatusAccepted)
		msg := messages.BuildGitRepository{
			Url:    submitModel.Url,
			Commit: submitModel.Commit,
		}

		h.buildQueue.Queue(&msg)

		return nil, nil
	}
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
