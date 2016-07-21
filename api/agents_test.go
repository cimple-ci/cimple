package api

import (
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_GetAgents(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Accept"][0] != "application/json" {
			t.Fatalf("Expected Accept header to be application/json - was %s", r.Header["Accept"])
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `[{"Id": "2bf80263-fa25-4df4-b132-aaa3bf62381d", "ConnectedDate": "2016-07-20T20:12:55.826456124+01:00", "Busy": true}]`)
	}))
	defer ts.Close()

	client, _ := NewApiClient()
	client.ServerUrl = ts.URL
	agents, err := client.GetAgents()

	if err != nil {
		t.Fatalf("Err %+v getting the agents", err)
	}

	if len(agents) != 1 {
		t.Fatalf("Expect 1 agent returned - was %d", len(agents))
	}

	exp, _ := uuid.FromString("2bf80263-fa25-4df4-b132-aaa3bf62381d")
	if agents[0].Id != exp {
		t.Fatalf("Expected a different id - was %s", agents[0].Id)
	}

	expDate := time.Date(2016, time.July, 20, 20, 12, 55, 826456124, time.Local)
	if agents[0].ConnectedDate != expDate {
		t.Fatalf("Expected a different ConnectedDate - was %s expected %s", agents[0].ConnectedDate, expDate)
	}

	if !agents[0].Busy {
		t.Fatal("Expected agent to be busy")
	}
}
