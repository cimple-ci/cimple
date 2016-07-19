package server

import (
	"log"

	"github.com/lukesmith/cimple/chore"
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
