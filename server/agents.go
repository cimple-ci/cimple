package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lukesmith/cimple/web_application"
	"github.com/satori/go.uuid"
)

type agentsHandler struct {
	agentPool AgentPool
	logger    *log.Logger
}

type agentModel struct {
	Id            uuid.UUID
	ConnectedDate time.Time
	Busy          bool
}

func registerAgents(app *web_application.Application, agentPool AgentPool, logger *log.Logger) {
	handler := &agentsHandler{
		agentPool: agentPool,
		logger:    logger,
	}

	app.Handle("/agents", handler.getAgents).Methods("GET").Name("agents")
	app.WebSocket("/agents/connection", handler.channelAgents).Methods("GET").Name("agent_connection")
}

func (h *agentsHandler) getAgents(app *web_application.Application, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	agents := make([]*agentModel, 0)
	pooledAgents, _ := h.agentPool.GetAgents()
	for _, agent := range pooledAgents {
		m := &agentModel{
			Id:            agent.Id,
			ConnectedDate: agent.ConnectedDate,
			Busy:          agent.IsBusy(),
		}

		agents = append(agents, m)
	}

	return agents, nil
}

func (h *agentsHandler) channelAgents(app *web_application.Application, socket *websocket.Conn, w http.ResponseWriter, r *http.Request) error {
	agentId, err := uuid.FromString(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Id in request (%s) is not a valid uuid.", r.URL.Query().Get("id"))
		// TODO: Return relevant status code
		return err
	}

	conn := newWebsocketAgentConnection(socket, h.logger)
	agent := newAgent(agentId, conn, h.logger)

	h.agentPool.Join(agent)
	defer func() {
		h.agentPool.Leave(agent)
	}()

	agent.listen()

	return nil
}
