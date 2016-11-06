package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

type fakeBuildQueue struct {
	queued []BuildJob
}

func (bq *fakeBuildQueue) Queue(job BuildJob) {
	bq.queued = append(bq.queued, job)
}

func (bq *fakeBuildQueue) GetQueued() ([]BuildJob, error) {
	return bq.queued, nil
}

func Test_SubmitBuild(t *testing.T) {
	app, server := newWebApplication()
	buildQueue := &fakeBuildQueue{}
	buildQueue.queued = make([]BuildJob, 0)
	registerBuilds(app, nil, buildQueue, log.New(os.Stdout, "", 0))

	body := make(map[string]interface{})
	body["Url"] = "https://test.local"
	body["Commit"] = "master"

	buildsUrl := fmt.Sprintf("%s/builds", server.URL)
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

		job := buildQueue.queued[0].(*buildGitRepositoryJob)

		if job.Url != body["Url"] {
			t.Fatalf("Expected queued build to have Url %s - was %s", body["Url"], job.Url)
		}

		if job.Commit != body["Commit"] {
			t.Fatalf("Expected queued build to have commit %s - was %s", body["Commit"], job.Commit)
		}
	}
}

func Test_ListBuilds(t *testing.T) {
	app, server := newWebApplication()
	buildQueue := &fakeBuildQueue{}
	queuedItem := NewBuildGitRepositoryJob("https://test.git", "master")
	buildQueue.queued = []BuildJob{queuedItem}
	registerBuilds(app, nil, buildQueue, log.New(os.Stdout, "", 0))

	var reader io.Reader
	buildsUrl := fmt.Sprintf("%s/builds", server.URL)
	request, err := http.NewRequest("GET", buildsUrl, reader)
	request.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(request)

	assert := assert.New(t)

	if assert.Nil(err) {
		assert.Equal(200, res.StatusCode, "OK expected")
		assert.Equal("application/json", res.Header["Content-Type"][0])

		var m []map[string]interface{}
		json.NewDecoder(res.Body).Decode(&m)
		assert.Equal(queuedItem.Id().String(), m[0]["id"])

		submissionDate, _ := time.Parse(time.RFC3339, m[0]["submission_date"].(string))
		assert.Equal(queuedItem.SubmissionDate().UTC(), submissionDate, "they should equal")

		assert.Equal("http://cimple.test/builds/"+queuedItem.Id().String(), m[0]["build_url"])
	}
}
