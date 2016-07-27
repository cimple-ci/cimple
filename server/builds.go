package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lukesmith/cimple/database"
	"github.com/lukesmith/cimple/web_application"
	"github.com/satori/go.uuid"
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

type buildsItemModel struct {
	Id             uuid.UUID `json:"id"`
	SubmissionDate time.Time `json:"submission_date"`
	BuildUrl       string    `json:"build_url"`
}

type submitBuildModel struct {
	Url    string `json:"url"`
	Commit string `json:"commit"`
}

func registerBuilds(app *web_application.Application, db database.CimpleDatabase, buildQueue BuildQueue, logger *log.Logger) {
	handler := &buildsHandler{
		db:         db,
		buildQueue: buildQueue,
		logger:     logger,
	}

	app.Handle("/builds/{key}", handler.getDetails).Methods("GET").Name("build")
	app.Handle("/builds", handler.listBuilds).Methods("GET").Name("listBuilds")
	app.Handle("/builds", handler.submitBuild).Methods("POST").Name("submitBuild")
}

func (h *buildsHandler) listBuilds(app *web_application.Application, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	pendingBuilds, err := h.buildQueue.GetQueued()
	if err != nil {
		return nil, err
	}

	builds := make([]*buildsItemModel, 0)

	for _, build := range pendingBuilds {
		//projectUrl, _ := app.Router.Get("project").URL("key", "project")
		buildUrl, _ := app.Router.Get("build").URL("key", build.Id().String())

		buildModel := &buildsItemModel{
			Id:             build.Id(),
			SubmissionDate: build.SubmissionDate(),
			BuildUrl:       buildUrl.String(),
		}

		builds = append(builds, buildModel)
	}

	return builds, nil
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
		job := NewBuildGitRepositoryJob(submitModel.Url, submitModel.Commit)

		h.buildQueue.Queue(job)

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
