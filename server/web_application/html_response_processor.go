package web_application

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/unrolled/render"
)

type HtmlResponseProcessor struct {
	render *render.Render
}

func NewHtmlResponseProcessor(render *render.Render) *HtmlResponseProcessor {
	return &HtmlResponseProcessor{
		render: render,
	}
}

func (p *HtmlResponseProcessor) CanProcess(mediaRange string) bool {
	return strings.EqualFold(mediaRange, render.ContentHTML) ||
		strings.EqualFold(mediaRange, render.ContentXHTML)
}

func (p *HtmlResponseProcessor) Process(w http.ResponseWriter, model interface{}, errorHandler func(w http.ResponseWriter, err error)) {
	templateName := reflect.TypeOf(model).Name()

	if strings.HasSuffix(templateName, "Model") {
		templateName = strings.Split(templateName, "Model")[0]

		p.render.HTML(w, http.StatusOK, templateName, model)
	} else {
		errorHandler(w, fmt.Errorf("Unable to find view template for %s. Does not have the Model suffix.", templateName))
	}
}
