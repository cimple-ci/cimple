package api

import (
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
