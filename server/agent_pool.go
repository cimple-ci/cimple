package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lukesmith/cimple/chore"
	"github.com/satori/go.uuid"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type agentpool struct {
	join       chan *Agent
	leave      chan *Agent
	agents     map[*Agent]bool
	logger     *log.Logger
	workerpool *chore.WorkPool
}

func (a *agentpool) GetAgents() ([]*Agent, error) {
	r := []*Agent{}
	for agent := range a.agents {
		r = append(r, agent)
	}
	return r, nil
}

func (a *agentpool) run() {
	for {
		select {
		case agent := <-a.join:
			a.agents[agent] = true
			a.workerpool.AddWorker(agent)
			a.logger.Printf("Agent %s joined", agent.Id)
		case agent := <-a.leave:
			delete(a.agents, agent)
			a.workerpool.RemoveWorker(agent)
			a.logger.Printf("Agent %s left", agent.Id)
		}
	}
}

func newAgentPool(logger *log.Logger) *agentpool {
	pool := &agentpool{
		join:       make(chan *Agent),
		leave:      make(chan *Agent),
		agents:     make(map[*Agent]bool),
		logger:     logger,
		workerpool: chore.NewWorkPool(),
	}

	return pool
}

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (s *agentpool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Print(err)
	}

	agentId, err := uuid.FromString(r.URL.Query().Get("id"))
	if err != nil {
		s.logger.Printf("Id in request (%s) is not a valid uuid.", r.URL.Query().Get("id"))
		// TODO: Return relevant status code
		return
	}

	conn := newWebsocketAgentConnection(socket, s.logger)
	agent := newAgent(agentId, conn, s.logger)

	s.join <- agent
	defer func() {
		s.leave <- agent
	}()

	agent.listen()
}
