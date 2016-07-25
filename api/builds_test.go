package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_SubmitBuild(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Accept"][0] != "application/json" {
			t.Fatalf("Expected Accept header to be application/json - was %s", r.Header["Accept"])
		}

		var m map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&m)
		r.Body.Close()

		if err != nil {
			t.Fatalf("%+v", err)
		}

		if m["Url"] != "https://test.local" {
			t.Fatalf("Unexpected Url value - %s", m["Url"])
		}

		if m["Commit"] != "master" {
			t.Fatalf("Unexpected Url value - %s", m["Commit"])
		}

		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	client, _ := NewApiClient()
	client.ServerUrl = ts.URL

	submissionOptions := &BuildSubmissionOptions{
		Url:    "https://test.local",
		Commit: "master",
	}

	err := client.SubmitBuild(submissionOptions)
	if err != nil {
		t.Fatalf("Err %+v submitting a build", err)
	}
}
