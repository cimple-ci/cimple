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
	Id        uuid.UUID
	conn      AgentConnection
	logger    *log.Logger
	busy      bool
	available chan bool
	router    *messages.Router
}

func (worker *Agent) CanPerform(c *chore.Chore) bool {
	return !worker.busy
}

func (worker *Agent) Perform(c *chore.Chore) error {
	worker.busy = true
	worker.send(c.Job)

	<-worker.available
	worker.logger.Printf("Agent now available")
	worker.busy = false

	return nil
}

func newAgent(agentId uuid.UUID, conn AgentConnection, logger *log.Logger) *Agent {
	agent := &Agent{
		Id:        agentId,
		conn:      conn,
		logger:    logger,
		available: make(chan bool),
		router:    messages.NewRouter(),
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
	c.logger.Printf("ServerAgent:%s - Sending %s:%s", c, name, env.Id)

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

	agent.router.On(messages.RegisterAgentMessage{}, func(m interface{}) {
		msg := m.(messages.RegisterAgentMessage)
		agent.send(&messages.ConfirmationMessage{
			ConfirmedId: msg.Id,
			Text:        "Thankyou from server"})
	})

	agent.router.On(messages.BuildComplete{}, func(m interface{}) {
		agent.available <- true
	})

	for {
		if msg, err := agent.read(); err == nil {
			name := reflect.TypeOf(msg.Body).Name()
			agent.logger.Printf("ServerAgent:%s - Received %s:%s", agent, name, msg.Id)
			agent.router.Route(msg.Body)
		} else {
			log.Printf("uhoh - %s", err)
			break
		}
	}
}
