package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type agentpool struct {
	join   chan *agent
	leave  chan *agent
	agents map[*agent]bool
	logger *log.Logger
}

func (a *agentpool) run() {
	for {
		select {
		case agent := <-a.join:
			a.agents[agent] = true
			a.logger.Print("Agent joined")
		case agent := <-a.leave:
			delete(a.agents, agent)
			a.logger.Print("Agent left")
		}
	}
}

func newAgentPool(logger *log.Logger) *agentpool {
	return &agentpool{
		join:   make(chan *agent),
		leave:  make(chan *agent),
		agents: make(map[*agent]bool),
		logger: logger,
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (s *agentpool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Print(err)
	}

	agent := &agent{
		socket: socket,
		pool:   s,
	}

	s.join <- agent
	defer func() {
		s.leave <- agent
	}()
	agent.read(s.logger)
}
