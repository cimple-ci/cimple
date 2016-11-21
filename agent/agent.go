package agent

import (
	"io"
	"log"
	"time"

	"github.com/kardianos/osext"
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
	SyslogUrl  string
}

func DefaultConfig() (*Config, error) {
	c := &Config{}
	return c, nil
}

type Agent struct {
	Id        uuid.UUID
	config    *Config
	logger    *log.Logger
	conn      *serverConnection
	router    *messages.Router
	connected bool
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

func NewAgent(config *Config, logger *log.Logger) (*Agent, error) {
	a := &Agent{
		Id:     uuid.NewV4(),
		config: config,
		logger: logger,
		router: messages.NewRouter(),
	}

	return a, nil
}

func (agent *Agent) Listen() {
	go func() {
		for {
			msg, err := agent.read()
			if err != nil {
				agent.logger.Printf("Error reading: %+v", err)
				break
			} else {
				name := reflect.TypeOf(msg.Body).Name()
				agent.logger.Printf("Received %s:%s", name, msg.Id)
				agent.router.Route(msg.Body)
			}
		}

		agent.logger.Printf("Stopped listening")
	}()
}

func (agent *Agent) Start() error {
	agent.logger.Printf("Starting agent %s", agent)

	agent.router.OnError(func(m interface{}) {
		agent.logger.Printf("Received an error routing %+v", m)
	})

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

		s, err := syslog.Dial("tcp", agent.config.SyslogUrl, syslog.LOG_INFO, "Runner", nil)
		if err != nil {
			agent.logger.Printf("Error connecting to syslog %+v", err)
		}
		defer s.Close()
		errWriter := io.MultiWriter(s)

		err = executeCimpleRun(pat, outWriter, errWriter)
		if err != nil {
			agent.logger.Printf("Err performing Cimple run %+v", err)
		}

		err = agent.send(&messages.BuildComplete{})
		if err != nil {
			agent.logger.Printf("Err sending build complete %+v", err)
		}
	})
	agent.router.On(messages.ConfirmationMessage{}, func(m interface{}) {
		msg := m.(messages.ConfirmationMessage)
		agent.logger.Printf("Confirmed %s", msg.Text)
	})

	conn, err := newWebsocketServerConnection(agent.config.ServerAddr, agent.config.ServerPort, agent.Id, agent.logger)
	if err != nil {
		return err
	}

	maintainConnection(agent, conn)

	for {
	}
}

func maintainConnection(agent *Agent, conn *serverConnection) error {
	agent.conn = conn

	go func() {
		for {
			select {
			case <-agent.conn.Connected:
				agent.logger.Printf("Agent connected")
				agent.Listen()
				agent.Register()
			case <-agent.conn.Disconnected:
				agent.logger.Print("Agent disconnected")
				reconnect(agent)
			}
		}
	}()

	err := agent.conn.Connect()
	if err != nil {
		reconnect(agent)
	}

	return nil
}

func reconnect(agent *Agent) {
	select {
	case <-time.After(time.Second * 1):
		conn, err := newWebsocketServerConnection(agent.config.ServerAddr, agent.config.ServerPort, agent.Id, agent.logger)
		if err != nil {
			agent.logger.Printf("Unable to connect %+v", err)
		}
		maintainConnection(agent, conn)
	}
}

func (agent Agent) Register() error {
	hostname, _ := os.Hostname()
	return agent.send(&messages.RegisterAgentMessage{
		Id:       agent.Id,
		Hostname: hostname,
	})
}

func executeCimpleRun(workingDir string, stdout io.Writer, stderr io.Writer) error {
	args := []string{"run", "--run-context", "server", "--journal-driver", "console", "--journal-format", "json"}
	filename, _ := osext.Executable()
	var cmd = exec.Command(filename, args...)
	cmd.Dir = workingDir
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
