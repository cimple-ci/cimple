package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
)

type ApiClient struct {
	ServerUrl string
}

type Agent struct {
	Id            uuid.UUID
	ConnectedDate time.Time
	Busy          bool
}

func NewApiClient() (*ApiClient, error) {
	return &ApiClient{}, nil
}

func (api *ApiClient) GetAgents() ([]Agent, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", api.ServerUrl+"/agents", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
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
