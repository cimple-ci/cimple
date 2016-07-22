package api

import (
	"bytes"
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

func (api *ApiClient) newHttpClient() *http.Client {
	client := &http.Client{}
	return client
}

func (api *ApiClient) newGetRequest(path string) (*http.Request, error) {
	req, err := http.NewRequest("GET", api.ServerUrl+"/"+path, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "cimple/cli")
	req.Header.Add("Accept", "application/json")

	return req, nil
}

func (api *ApiClient) newPostRequest(path string, body interface{}) (*http.Request, error) {
	reader := new(bytes.Buffer)
	json.NewEncoder(reader).Encode(body)
	req, err := http.NewRequest("POST", api.ServerUrl+"/"+path, reader)

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "cimple/cli")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}
