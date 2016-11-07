package agent

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lukesmith/cimple/messages"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type ServerConnection interface {
	SendMessage(env *messages.Envelope) error
	ReadMessage(env *messages.Envelope) error
	Close() error
}

type serverConnection struct {
	socket       *websocket.Conn
	encoder      *gob.Encoder
	decoder      *gob.Decoder
	logger       *log.Logger
	url          string
	agentId      uuid.UUID
	Connected    chan websocket.Conn
	Disconnected chan websocket.Conn
}

func newWebsocketServerConnection(addr string, port string, agentId uuid.UUID, logger *log.Logger) (*serverConnection, error) {
	url := fmt.Sprintf("ws://%s:%s/agents/connection?id=%s", addr, port, agentId)

	conn := &serverConnection{
		logger: logger,
		url:    url,
	}
	conn.encoder = gob.NewEncoder(conn)
	conn.decoder = gob.NewDecoder(conn)
	conn.Connected = make(chan websocket.Conn)
	conn.Disconnected = make(chan websocket.Conn)

	return conn, nil
}

func (a *serverConnection) ReadMessage(env *messages.Envelope) error {
	return a.decoder.Decode(env)
}

func (a *serverConnection) SendMessage(env *messages.Envelope) error {
	return a.encoder.Encode(env)
}

func (a *serverConnection) Write(p []byte) (n int, err error) {
	err = a.socket.WriteMessage(websocket.BinaryMessage, p)
	return len(p), err
}

func (a *serverConnection) Read(p []byte) (n int, err error) {
	_, message, err := a.socket.ReadMessage()
	if err != nil {
		return 0, err
	}

	for i := 0; i < len(message); i++ {
		p[i] = message[i]
	}

	return len(message), nil
}

func (a *serverConnection) Close() error {
	return a.socket.Close()
}

func (a *serverConnection) Connect() error {
	a.logger.Print("Attempting to connect to server")
	socket, _, err := websocket.DefaultDialer.Dial(a.url, nil)
	if err != nil {
		a.logger.Printf("Error connecting to server - %s", err)
		return err
	}

	a.socket = socket
	a.logger.Print("Connected")
	a.Connected <- *socket

	go a.PingServer()

	return nil
}

func (a *serverConnection) reconnect() error {
	return a.Connect()
}

func (a *serverConnection) PingServer() {
	a.logger.Print("Setting up server ping")
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	a.socket.SetPongHandler(func(appData string) error {
		a.logger.Print("Recieved Pong from server")
		return nil
	})

	for {
		end := false
		select {
		case <-ticker.C:
			a.logger.Print("Sending Ping")
			a.socket.SetWriteDeadline(time.Now().Add(pongWait))
			if err := a.socket.WriteMessage(websocket.PingMessage, []byte(a.agentId.String())); err != nil {
				a.logger.Println("Ping: ", err)
				end = true
			}
		}

		if end {
			ticker.Stop()
			break
		}
	}

	a.logger.Print("Stopping ping")
	a.socket.Close()
	a.Disconnected <- *a.socket
}
