package agent

import (
	"crypto/tls"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lukesmith/cimple/messages"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

type ServerConnection interface {
	SendMessage(env *messages.Envelope) error
	ReadMessage(env *messages.Envelope) error
	Close() error
}

type serverConnection struct {
	socket          *websocket.Conn
	encoder         *gob.Encoder
	decoder         *gob.Decoder
	logger          *log.Logger
	url             string
	agentId         uuid.UUID
	Connected       chan websocket.Conn
	Disconnected    chan websocket.Conn
	TLSClientConfig *tls.Config
}

func newWebsocketServerConnection(cfg *Config, agentId uuid.UUID, logger *log.Logger) (*serverConnection, error) {
	scheme := "ws"
	if cfg.EnableTLS {
		scheme = "wss"
	}
	url := fmt.Sprintf("%s://%s:%s/agents/connection?id=%s", scheme, cfg.ServerAddr, cfg.ServerPort, agentId)

	conn := &serverConnection{
		logger:          logger,
		url:             url,
		TLSClientConfig: cfg.TLSClientConfig,
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

	dialer := &websocket.Dialer{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: a.TLSClientConfig,
	}
	socket, _, err := dialer.Dial(a.url, nil)
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
