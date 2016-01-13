package web_application

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jchannon/negotiator"
	"github.com/unrolled/render"
)

type handler func(*Application, http.ResponseWriter, *http.Request) (interface{}, error)

type Application struct {
	render     *render.Render
	Router     *mux.Router
	negotiator *negotiator.Negotiator
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.ServeHTTP(w, r)
}

func NewApplication() *Application {
	router := mux.NewRouter()
	helpers := NewAppHelper(router)
	render := NewRenderer(helpers)
	neg := negotiator.New(NewHtmlResponseProcessor(render))

	return &Application{
		render:     render,
		Router:     router,
		negotiator: neg,
	}
}

func (app *Application) Handle(path string, handler handler) *mux.Route {
	return app.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		model, err := handler(app, w, r)

		if err != nil {
			GlobalErrorHandler(w, err)
		} else {
			app.negotiator.Negotiate(w, r, model, GlobalErrorHandler)
		}
	})
}

func (app *Application) Static(path string) *mux.Route {
	return app.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		file, err := ioutil.ReadFile("server/frontend/" + r.URL.Path)
		if err != nil {
			GlobalErrorHandler(w, err)
		} else {
			if strings.HasSuffix(r.URL.Path, ".js") {
				w.Header().Add("Content-Type", "application/javascript")
			}
			w.Write(file)
		}
	})
}

func GlobalErrorHandler(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func NewRenderer(helpers *AppHelpers) *render.Render {
	return render.New(render.Options{
		IsDevelopment: true,
		Layout:        "layout",
		Directory:     "server/frontend/templates",
		Funcs: []template.FuncMap{
			{"Url": helpers.URL},
		},
	})
}

type AppHelpers struct {
	router *mux.Router
}

func NewAppHelper(router *mux.Router) *AppHelpers {
	return &AppHelpers{
		router: router,
	}
}

func (ah *AppHelpers) URL(name string, pairs ...string) (string, error) {
	url, err := ah.router.Get(name).URL(pairs...)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
