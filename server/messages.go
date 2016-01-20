package server

import (
	"encoding/gob"
	"github.com/satori/go.uuid"
)

func init() {
	gob.Register(ConfirmationMessage{})
	gob.Register(RegisterAgentMessage{})
}

type Envelope struct {
	Id   uuid.UUID
	Body interface{}
}

type ConfirmationMessage struct {
	ConfirmedId uuid.UUID
	Text        string
}

type RegisterAgentMessage struct {
	Hostname string
}
