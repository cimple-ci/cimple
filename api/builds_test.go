package api

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func Test_ListBuilds(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Accept"][0] != "application/json" {
			t.Fatalf("Expected Accept header to be application/json - was %s", r.Header["Accept"])
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `[{"id": "2bf80263-fa25-4df4-b132-aaa3bf62381d", "submission_date": "2016-08-20T11:12:13.826456124+01:00"}]`)
	}))
	defer ts.Close()

	client, _ := NewApiClient()
	client.ServerUrl = ts.URL

	builds, err := client.ListBuilds()
	if err != nil {
		t.Fatalf("Err %+v getting the list of builds", err)
	}

	if len(builds) != 1 {
		t.Fatalf("Expect 1 build to be returned - was %d", len(builds))
	}

	exp, _ := uuid.FromString("2bf80263-fa25-4df4-b132-aaa3bf62381d")
	if builds[0].Id != exp {
		t.Fatalf("Expected a different id - was %s", builds[0].Id)
	}

	expSubmission := time.Date(2016, time.August, 20, 11, 12, 13, 826456124, time.Local)
	if builds[0].SubmissionDate != expSubmission {
		t.Fatalf("Expected a different SubmissionDate - was %s expected %s", builds[0].SubmissionDate, expSubmission)
	}
}
