package server

import (
	"github.com/satori/go.uuid"
	"time"

	"github.com/lukesmith/cimple/messages"
	"log"
	"reflect"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

type agent struct {
	id     uuid.UUID
	conn   AgentConnection
	logger *log.Logger
}

func newAgent(agentId uuid.UUID, conn AgentConnection, logger *log.Logger) *agent {
	agent := &agent{
		id:     agentId,
		conn:   conn,
		logger: logger,
	}
	return agent
}

func (a *agent) String() string {
	return a.id.String()
}

func (c *agent) send(msg interface{}) error {
	env := &messages.Envelope{
		Id:   uuid.NewV4(),
		Body: msg,
	}

	name := reflect.TypeOf(msg).Elem().Name()
	c.logger.Printf("Sending %s:%s to agent:%s", name, env.Id, c)

	return c.conn.SendMessage(env)
}

func (a *agent) read() (messages.Envelope, error) {
	var m messages.Envelope
	if err := a.conn.ReadMessage(&m); err == nil {
		return m, nil
	} else {
		return m, err
	}
}

func (agent *agent) listen(logger *log.Logger) {
	defer agent.conn.Close()

	for {
		if msg, err := agent.read(); err == nil {
			name := reflect.TypeOf(msg.Body).Name()
			logger.Printf("Received %s:%s from agent:%s", name, msg.Id, agent)
			agent.send(&messages.ConfirmationMessage{
				ConfirmedId: msg.Id,
				Text:        "Thankyou from server"})
		} else {
			log.Printf("uhoh - %s", err)
			break
		}
	}
}
