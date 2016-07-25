package api

import (
	"fmt"
)

type BuildSubmissionOptions struct {
	Url    string
	Commit string
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
