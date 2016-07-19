package server

import (
	"net/http"
	"log"

	"github.com/gorilla/websocket"
	"github.com/lukesmith/cimple/web_application"
	"github.com/satori/go.uuid"
)

type agentsHandler struct {
	agentPool *agentpool
}

func registerAgents(app *web_application.Application, agentPool *agentpool) {
	handler := &agentsHandler{
		agentPool: agentPool,
	}

	app.WebSocket("/agents/connection", handler.channelAgents).Methods("GET").Name("agent_connection")
}

func (h *agentsHandler) getAgents(app *web_application.Application, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, nil
}


func (h *agentsHandler) channelAgents(app *web_application.Application, socket *websocket.Conn, w http.ResponseWriter, r *http.Request) error {
	agentId, err := uuid.FromString(r.URL.Query().Get("id"))
	if err != nil {
		log.Printf("Id in request (%s) is not a valid uuid.", r.URL.Query().Get("id"))
		// TODO: Return relevant status code
		return err
	}

	conn := newWebsocketAgentConnection(socket, h.agentPool.logger)
	agent := newAgent(agentId, conn, h.agentPool.logger)

	h.agentPool.join <- agent
	defer func() {
		h.agentPool.leave <- agent
	}()

	agent.listen()

	return nil
}
