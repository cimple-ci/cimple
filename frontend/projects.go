package frontend

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lukesmith/cimple/database"
	"github.com/lukesmith/cimple/web_application"
)

type projectsHandler struct {
	db database.CimpleDatabase
}

type projectModel struct {
	Name       string        `json:"name"`
	ProjectUrl string        `json:"project_url"`
	BuildCount int           `json:"build_count"`
	Builds     []*buildModel `json:"builds"`
}

func registerProjects(app *web_application.Application, db database.CimpleDatabase) {
	handler := &projectsHandler{
		db: db,
	}

	app.Handle("/projects", handler.getIndex).Methods("GET").Name("projects")
	app.Handle("/projects/{key}", handler.getDetails).Methods("GET").Name("project")
}

func (h *projectsHandler) getIndex(app *web_application.Application, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	projects := h.db.GetProjects()

	model := []*projectModel{}

	for _, proj := range projects {
		url, _ := app.Router.Get("project").URL("key", proj.Name)
		model = append(model, &projectModel{Name: proj.Name, ProjectUrl: url.Path, BuildCount: proj.BuildCount})
	}

	return model, nil
}

func (h *projectsHandler) getDetails(app *web_application.Application, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)

	p, err := h.db.GetProject(params["key"])
	if err != nil {
		return nil, err
	}

	builds, err := h.db.GetBuilds(p.Name)
	if err != nil {
		return nil, err
	}

	buildModels := []*buildModel{}

	for _, build := range builds {
		url, _ := app.Router.Get("project").URL("key", p.Name)
		buildUrl, _ := app.Router.Get("build").URL("project_key", p.Name, "key", build.Id)
		buildModels = append(buildModels, &buildModel{Id: p.Name, ProjectUrl: url.Path, BuildUrl: buildUrl.Path})
	}

	url, _ := app.Router.Get("project").URL("key", p.Name)
	return projectModel{Name: p.Name, ProjectUrl: url.Path, BuildCount: p.BuildCount, Builds: buildModels}, nil
}
