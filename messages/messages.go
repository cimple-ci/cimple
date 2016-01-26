package messages

import (
	"encoding/gob"
	"github.com/satori/go.uuid"
)

func init() {
	gob.Register(ConfirmationMessage{})
	gob.Register(RegisterAgentMessage{})
	gob.Register(BuildGitRepository{})
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

type BuildGitRepository struct {
	Url    string
	Commit string
}
