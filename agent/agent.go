package agent

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lukesmith/cimple/logging"
	"github.com/satori/go.uuid"
	"os"
	"reflect"
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
	Id     uuid.UUID
	config *Config
	logger *log.Logger
	conn   *serverConnection
}

func (a *Agent) String() string {
	return a.Id.String()
}

func (c *Agent) send(msg interface{}) error {
	env := &Envelope{
		Id:   uuid.NewV4(),
		Body: msg,
	}

	name := reflect.TypeOf(msg).Elem().Name()
	c.logger.Printf("Sending %s:%s", name, env.Id)

	return c.conn.SendMessage(env)
}

func (a *Agent) read() (Envelope, error) {
	var m Envelope
	if err := a.conn.ReadMessage(&m); err == nil {
		return m, nil
	} else {
		return m, err
	}
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
	agent.conn = newWebsocketServerConnection(c, agent.logger)

	defer agent.conn.Close()

	go func() {
		defer agent.conn.Close()
		for {
			msg, err := agent.read()
			if err != nil {
				agent.logger.Printf("Err reading: %+v", err)
			} else {
				name := reflect.TypeOf(msg.Body).Name()
				agent.logger.Printf("Received %s:%s", name, msg.Id)
			}
		}
	}()

	agent.Register()
	agent.conn.PingServer(agent.Id)

	return nil
}

func (agent Agent) Register() error {
	hostname, _ := os.Hostname()
	return agent.send(&RegisterAgentMessage{
		Hostname: hostname,
	})
}
