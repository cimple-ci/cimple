package frontend

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lukesmith/cimple/server/web_application"
	"github.com/lukesmith/cimple/database"
)

type projectsHandler struct {
	db database.CimpleDatabase
}

type projectModel struct {
	Name       string `json:"name"`
	ProjectUrl string `json:"project_url"`
}

func RegisterProjects(app *web_application.Application, db database.CimpleDatabase) {
	handler := &projectsHandler{
		db: db,
	}

	app.Handle("/projects", handler.getIndex).Methods("GET").Name("projects")
	app.Handle("/projects/{key}", handler.getDetails).Methods("GET").Name("project")
}

func(h *projectsHandler) getIndex(app *web_application.Application, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	projects := h.db.GetProjects()

	model := []*projectModel{}

	for _, proj := range projects {
		url, _ := app.Router.Get("project").URL("key", proj.Name)
		model = append(model, &projectModel{Name: proj.Name, ProjectUrl: url.Path})
	}

	return model, nil
}

func(h *projectsHandler) getDetails(app *web_application.Application, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	return params["key"], nil
}
