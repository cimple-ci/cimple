package server

import (
	"github.com/satori/go.uuid"
	"time"

	"github.com/lukesmith/cimple/chore"
	"github.com/lukesmith/cimple/messages"
	"log"
	"reflect"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

type Agent struct {
	Id     uuid.UUID
	conn   AgentConnection
	logger *log.Logger
	busy   bool
}

func (worker *Agent) CanPerform(c *chore.Chore) bool {
	return true
}

func (worker *Agent) Perform(c *chore.Chore) error {
	worker.busy = true
	defer func() {
		worker.busy = false
	}()
	log.Printf("Performing chore %d", c.ID)
	time.Sleep(7000 * time.Millisecond)
	log.Printf("Performed chore %d", c.ID)
	return nil
}

func newAgent(agentId uuid.UUID, conn AgentConnection, logger *log.Logger) *Agent {
	agent := &Agent{
		Id:     agentId,
		conn:   conn,
		logger: logger,
	}
	return agent
}

func (a *Agent) String() string {
	return a.Id.String()
}

func (c *Agent) send(msg interface{}) error {
	env := &messages.Envelope{
		Id:   uuid.NewV4(),
		Body: msg,
	}

	name := reflect.TypeOf(msg).Elem().Name()
	c.logger.Printf("Agent:%s - Sending %s:%s", c, name, env.Id)

	return c.conn.SendMessage(env)
}

func (a *Agent) read() (messages.Envelope, error) {
	var m messages.Envelope
	if err := a.conn.ReadMessage(&m); err == nil {
		return m, nil
	} else {
		return m, err
	}
}

func (agent *Agent) listen() {
	defer agent.conn.Close()

	for {
		if msg, err := agent.read(); err == nil {
			name := reflect.TypeOf(msg.Body).Name()
			agent.logger.Printf("Agent:%s - Received %s:%s", agent, name, msg.Id)
			agent.send(&messages.ConfirmationMessage{
				ConfirmedId: msg.Id,
				Text:        "Thankyou from server"})
		} else {
			log.Printf("uhoh - %s", err)
			break
		}
	}
}
