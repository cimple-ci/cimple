package web_application

import (
	"errors"
	"github.com/unrolled/render"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCanProcess(t *testing.T) {
	p := &HtmlResponseProcessor{}

	if !p.CanProcess(render.ContentHTML) {
		t.Fatalf("Expected to be able to process %s", render.ContentHTML)
	}

	if !p.CanProcess(render.ContentXHTML) {
		t.Fatalf("Expected to be able to process %s", render.ContentXHTML)
	}

	if p.CanProcess(render.ContentText) {
		t.Fatalf("Expected not to be able to process %s", render.ContentText)
	}
}

func TestProcess_ModelTypeDoesNotEndWithModel(t *testing.T) {
	p := &HtmlResponseProcessor{}

	recorder := httptest.NewRecorder()
	model := "hello"
	errorHappened := false

	p.Process(recorder, model, func(w http.ResponseWriter, err error) {
		errorHappened = true
	})

	if !errorHappened {
		t.Fatal("Expected an error to occur")
	}
}

func TestProcess_ModelTypeEndsWithModel(t *testing.T) {
	ren := render.New(render.Options{
		Asset: func(name string) ([]byte, error) {
			if name == "templates/fake.tmpl" {
				return []byte("fake {{.Name}}"), nil
			}

			return []byte{}, errors.New("Unknown asset")
		},
		AssetNames: func() []string {
			return []string{"templates/fake.tmpl"}
		},
	})

	p := NewHtmlResponseProcessor(ren)
	recorder := httptest.NewRecorder()
	model := fakeModel{
		Name: "John",
	}

	p.Process(recorder, model, func(w http.ResponseWriter, err error) {
		t.Fatalf("Did not expect an error to occur - %#v", err)
	})

	if recorder.Code != 200 {
		t.Fatalf("Expected recorded code to be 200, was %d", recorder.Code)
	}

	content_type := recorder.HeaderMap.Get(render.ContentType)
	if content_type != render.ContentHTML+"; charset=UTF-8" {
		t.Fatalf("Expected %s to be set to %s, but was %s", render.ContentType, render.ContentHTML, content_type)
	}

	if recorder.Body.String() != "fake John" {
		t.Fatalf("Expected body content to be fake John, but was %s", recorder.Body.String())
	}
}

type fakeModel struct {
	Name string
}
