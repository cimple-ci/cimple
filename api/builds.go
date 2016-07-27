package api

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

type BuildSubmissionOptions struct {
	Url    string
	Commit string
}

type Build struct {
	Id             uuid.UUID `json:"id"`
	SubmissionDate time.Time `json:"submission_date"`
}

func (api *ApiClient) SubmitBuild(options *BuildSubmissionOptions) error {
	client := api.newHttpClient()
	req, err := api.newPostRequest("builds", options)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return fmt.Errorf("Non accepted response %d", resp.StatusCode)
	}

	return nil
}

func (api *ApiClient) ListBuilds() ([]Build, error) {
	client := api.newHttpClient()
	req, err := api.newGetRequest("builds")
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var record []Build
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		return nil, err
	}

	return record, nil
}
