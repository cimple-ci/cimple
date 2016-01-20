package server

import (
	"encoding/gob"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"time"
	"github.com/lukesmith/cimple/messages"
)

type AgentConnection interface {
	SendMessage(env *messages.Envelope) error
	ReadMessage(env *messages.Envelope) error
	Close() error
}

type agentConnection struct {
	socket  *websocket.Conn
	encoder *gob.Encoder
	decoder *gob.Decoder
	logger  *log.Logger
}

func newWebsocketAgentConnection(socket *websocket.Conn, logger *log.Logger) *agentConnection {
	conn := &agentConnection{
		socket: socket,
		logger: logger,
	}
	conn.encoder = gob.NewEncoder(conn)
	conn.decoder = gob.NewDecoder(conn)

	conn.socket.SetPingHandler(func(appData string) error {
		logger.Printf("Received Ping from agent:%s", appData)
		logger.Printf("Sending Pong to agent:%s", appData)
		err := conn.socket.WriteControl(websocket.PongMessage, []byte(""), time.Now().Add(writeWait))
		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})

	return conn
}

func (a *agentConnection) ReadMessage(env *messages.Envelope) error {
	return a.decoder.Decode(env)
}

func (a *agentConnection) SendMessage(env *messages.Envelope) error {
	return a.encoder.Encode(env)
}

func (a *agentConnection) Write(p []byte) (n int, err error) {
	err = a.socket.WriteMessage(websocket.BinaryMessage, p)
	return len(p), err
}

func (a *agentConnection) Read(p []byte) (n int, err error) {
	_, message, err := a.socket.ReadMessage()
	if err != nil {
		return 0, err
	}

	for i := 0; i < len(message); i++ {
		p[i] = message[i]
	}

	return len(message), nil
}

func (a *agentConnection) Close() error {
	return a.socket.Close()
}
