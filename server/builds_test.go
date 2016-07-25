package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lukesmith/cimple/messages"
	"log"
	"net/http"
	"os"
	"testing"
)

var (
	buildsUrl string
)

func init() {
	buildsUrl = fmt.Sprintf("%s/builds", server.URL)
}

type fakeBuildQueue struct {
	queued []*messages.BuildGitRepository
}

func (bq *fakeBuildQueue) Queue(job interface{}) {
	a := job.(*messages.BuildGitRepository)
	bq.queued = append(bq.queued, a)
}

func Test_SubmitBuild(t *testing.T) {
	buildQueue := &fakeBuildQueue{}
	buildQueue.queued = make([]*messages.BuildGitRepository, 0)
	registerBuilds(app, nil, buildQueue, log.New(os.Stdout, "", 0))

	body := make(map[string]interface{})
	body["Url"] = "https://test.local"
	body["Commit"] = "master"

	reader := new(bytes.Buffer)
	json.NewEncoder(reader).Encode(body)
	request, err := http.NewRequest("POST", buildsUrl, reader)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	} else {
		if res.StatusCode != 202 {
			t.Errorf("Accepted expected: %d", res.StatusCode)
		}

		if len(buildQueue.queued) != 1 {
			t.Fatalf("Expected a build to have been queued")
		}

		job := buildQueue.queued[0]

		if job.Url != body["Url"] {
			t.Fatalf("Expected queued build to have Url %s - was %s", body["Url"], job.Url)
		}

		if job.Commit != body["Commit"] {
			t.Fatalf("Expected queued build to have commit %s - was %s", body["Commit"], job.Commit)
		}
	}
}
