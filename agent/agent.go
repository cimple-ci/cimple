package agent

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lukesmith/cimple/logging"
	"github.com/satori/go.uuid"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 1) / 10
)

type Config struct {
	ServerPort string
}

func DefaultConfig() (*Config, error) {
	c := &Config{}
	return c, nil
}

type Agent struct {
	Id         uuid.UUID
	config     *Config
	logger     *log.Logger
	serverConn *websocket.Conn
}

func (a *Agent) String() string {
	return a.Id.String()
}

func NewAgent(config *Config, logger io.Writer) (*Agent, error) {
	a := &Agent{
		Id:     uuid.NewV4(),
		config: config,
		logger: logging.CreateLogger("Agent", logger),
	}

	return a, nil
}

func (agent *Agent) Start() error {
	agent.logger.Printf("Starting agent %s", agent)

	url := fmt.Sprintf("ws://localhost:%s/agents?id=%s", agent.config.ServerPort, agent.Id)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	agent.serverConn = c

	go func() {
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				agent.logger.Print(err)
			}
			agent.logger.Printf("agent recv: %s", message)
		}
	}()

	agent.Register()
	agent.pingServer()

	return nil
}

func (agent *Agent) pingServer() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	ws := agent.serverConn

	ws.SetPongHandler(func(appData string) error {
		agent.logger.Printf("Pong: Agent rcv pong from server %s", appData)
		return nil
	})

	for {
		select {
		case <-ticker.C:
			agent.logger.Print("Ping: Sending")
			ws.SetWriteDeadline(time.Now().Add(pongWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte("From agent")); err != nil {
				agent.logger.Println("Ping: ", err)
			}
		}
	}
}

func (agent Agent) Register() error {
	return agent.serverConn.WriteMessage(websocket.TextMessage, []byte("Registering"))
}
