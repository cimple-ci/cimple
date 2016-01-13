package web_application

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jchannon/negotiator"
	"github.com/unrolled/render"
	"path/filepath"
)

type handler func(*Application, http.ResponseWriter, *http.Request) (interface{}, error)

type ApplicationOptions struct {
	ViewsDirectory string
}

type Application struct {
	render     *render.Render
	Router     *mux.Router
	negotiator *negotiator.Negotiator
	options    *ApplicationOptions
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.ServeHTTP(w, r)
}

func NewApplication(options *ApplicationOptions) *Application {
	router := mux.NewRouter()
	helpers := NewAppHelper(router)
	render := NewRenderer(options, helpers)
	neg := negotiator.New(NewHtmlResponseProcessor(render))

	return &Application{
		render:     render,
		Router:     router,
		negotiator: neg,
		options:    options,
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
		file, err := ioutil.ReadFile(filepath.Join(app.options.ViewsDirectory, r.URL.Path))
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

func NewRenderer(options *ApplicationOptions, helpers *AppHelpers) *render.Render {
	return render.New(render.Options{
		IsDevelopment: true,
		Layout:        "layout",
		Directory:     options.ViewsDirectory,
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
