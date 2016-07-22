package api

import (
	"encoding/json"
)

type BuildSubmissionOptions struct {
	Url    string
	Commit string
}

func (api *ApiClient) SubmitBuild(options BuildSubmissionOptions) error {
	client := api.newHttpClient()
	req, err := api.newPostRequest("builds", options)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var record []Agent
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		return err
	}

	return nil
}
