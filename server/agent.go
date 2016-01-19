package server

import (
	"github.com/satori/go.uuid"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

type agent struct {
	id     uuid.UUID
	socket *websocket.Conn
	pool   *agentpool
}

func (c *agent) send(msg []byte) error {
	return c.socket.WriteMessage(websocket.TextMessage, msg)
}

func (c *agent) read(logger *log.Logger) {
	defer c.socket.Close()
	c.socket.SetPingHandler(func(appData string) error {
		logger.Printf("Ping: Recieved - %s from %s", appData, c.id)
		err := c.socket.WriteControl(websocket.PongMessage, []byte("message"), time.Now().Add(writeWait))
		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})

	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			logger.Printf("Server rcv: %s", msg)
			c.send([]byte("Thankyou from server"))
		} else {
			log.Printf("uhoh - %s", err)
			break
		}
	}
}
