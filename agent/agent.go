package agent

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lukesmith/cimple/logging"
	"github.com/lukesmith/cimple/messages"
	"github.com/lukesmith/cimple/vcs/git"
	"github.com/lukesmith/syslog"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"os/exec"
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
	ServerAddr string
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
	router *messages.Router
}

func (a *Agent) String() string {
	return a.Id.String()
}

func (c *Agent) send(msg interface{}) error {
	env := &messages.Envelope{
		Id:   uuid.NewV4(),
		Body: msg,
	}

	name := reflect.TypeOf(msg).Elem().Name()
	c.logger.Printf("Sending %s:%s", name, env.Id)

	return c.conn.SendMessage(env)
}

func (a *Agent) read() (messages.Envelope, error) {
	var m messages.Envelope
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
		router: messages.NewRouter(),
	}

	return a, nil
}

func (agent *Agent) Start() error {
	agent.logger.Printf("Starting agent %s", agent)

	url := fmt.Sprintf("ws://%s:%s/agents?id=%s", agent.config.ServerAddr, agent.config.ServerPort, agent.Id)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	agent.conn = newWebsocketServerConnection(c, agent.logger)

	syslogUrl := fmt.Sprintf("%s:1514", agent.config.ServerAddr)

	s, err := syslog.Dial("tcp", syslogUrl, syslog.LOG_INFO, "Agent", nil)
	if err != nil {
		agent.logger.Printf("Failed to dial syslog on server = %+v", err)
	}
	defer s.Close()

	sOut, err := syslog.Dial("tcp", syslogUrl, syslog.LOG_INFO, "Runner", nil)
	if err != nil {
		agent.logger.Printf("Failed to dial syslog on server = %+v", err)
	}
	defer sOut.Close()

	sErr, err := syslog.Dial("tcp", syslogUrl, syslog.LOG_DEBUG, "Runner", nil)
	if err != nil {
		agent.logger.Printf("Failed to dial syslog on server = %+v", err)
	}
	defer sErr.Close()

	agent.logger = logging.CreateLogger("Agent", s)

	agent.router.On(messages.BuildGitRepository{}, func(m interface{}) {
		msg := m.(messages.BuildGitRepository)
		agent.logger.Printf("Building git repo:%s", msg.Url)

		pat, err := ioutil.TempDir("", "")
		if err != nil {
			agent.logger.Printf("Err %+v", err)
		}

		agent.logger.Printf("Creating dir %s", pat)
		err = os.MkdirAll(pat, 0755)
		if err != nil {
			agent.logger.Printf("Err %+v", err)
		}

		cloneOptions := git.NewCloneOptions(msg.Url, pat)
		err = git.Clone(cloneOptions, os.Stdout)
		if err != nil {
			agent.logger.Printf("Err during clone %+v", err)
		}

		checkoutOptions := git.NewCheckoutOptions(pat, msg.Commit)
		err = git.Checkout(checkoutOptions, os.Stdout)
		if err != nil {
			agent.logger.Printf("Err during checkout %+v", err)
		}

		outWriter := io.MultiWriter(os.Stdout)
		errWriter := io.MultiWriter(os.Stderr)

		// TODO: forward stdout as journal messages, stderr as syslog to server
		err = executeCimpleRun(pat, outWriter, errWriter, syslogUrl)
		if err != nil {
			agent.logger.Printf("Err performing Cimple run %+v", err)
		}
	})
	agent.router.On(messages.ConfirmationMessage{}, func(m interface{}) {
		msg := m.(messages.ConfirmationMessage)
		agent.logger.Printf("Confirmed %s", msg.Text)
	})

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
				agent.router.Route(msg.Body)
			}
		}
	}()

	agent.Register()
	agent.conn.PingServer(agent.Id)

	return nil
}

func (agent Agent) Register() error {
	hostname, _ := os.Hostname()
	return agent.send(&messages.RegisterAgentMessage{
		Hostname: hostname,
	})
}

func executeCimpleRun(workingDir string, stdout io.Writer, stderr io.Writer, syslogUrl string) error {
	args := []string{"run", "--syslog", syslogUrl}
	var cmd = exec.Command("cimple", args...)
	cmd.Dir = workingDir
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
