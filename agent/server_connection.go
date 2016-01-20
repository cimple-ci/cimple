package agent

import (
	"encoding/gob"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

type ServerConnection interface {
	SendMessage(env *Envelope) error
	ReadMessage(env *Envelope) error
	Close() error
}

type serverConnection struct {
	socket  *websocket.Conn
	encoder *gob.Encoder
	decoder *gob.Decoder
	logger  *log.Logger
}

func newWebsocketServerConnection(socket *websocket.Conn, logger *log.Logger) *serverConnection {
	conn := &serverConnection{
		socket: socket,
		logger: logger,
	}
	conn.encoder = gob.NewEncoder(conn)
	conn.decoder = gob.NewDecoder(conn)

	return conn
}

func (a *serverConnection) ReadMessage(env *Envelope) error {
	return a.decoder.Decode(env)
}

func (a *serverConnection) SendMessage(env *Envelope) error {
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

func (a *serverConnection) PingServer(agentId uuid.UUID) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	a.socket.SetPongHandler(func(appData string) error {
		a.logger.Print("Recieved Pong from server")
		return nil
	})

	for {
		select {
		case <-ticker.C:
			a.logger.Print("Sending Ping")
			a.socket.SetWriteDeadline(time.Now().Add(pongWait))
			if err := a.socket.WriteMessage(websocket.PingMessage, []byte(agentId.String())); err != nil {
				a.logger.Println("Ping: ", err)
			}
		}
	}
}
