package api

import (
	"encoding/json"
)

func (api *ApiClient) GetAgents() ([]Agent, error) {
	client := api.newHttpClient()
	req, err := api.newGetRequest("agents")
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var record []Agent
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		return nil, err
	}

	return record, nil
}
