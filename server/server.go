package server

import (
	"io"
	"log"
	"net"
	"net/http"

	"github.com/lukesmith/cimple/database"
	"github.com/lukesmith/cimple/frontend"
	"github.com/lukesmith/cimple/logging"
	"github.com/mcuadros/go-syslog"
)

type Config struct {
	Addr       string
	SyslogAddr string
}

func DefaultConfig() *Config {
	return &Config{
		Addr:       ":0",
		SyslogAddr: ":1514",
	}
}

type Server struct {
	config *Config
	logger *log.Logger
}

func NewServer(config *Config, logger io.Writer) (*Server, error) {
	s := &Server{
		config: config,
		logger: logging.CreateLogger("Server", logger),
	}

	return s, nil
}

func (server *Server) Start() error {
	agents := newAgentPool(server.logger)

	db := database.NewDatabase("./.cimple")
	app := frontend.NewFrontend(db)
	hooks := NewHooks(agents)

	http.Handle("/", app)
	http.Handle("/agents", agents)
	http.Handle("/hooks", hooks)

	go syslogEndpoint(server)

	go agents.run()

	s := &http.Server{
		Addr: server.config.Addr,
	}

	addr := server.config.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	server.logger.Printf("Cimple server available on %s", ln.Addr())
	err = s.Serve(ln)
	if err != nil {
		server.logger.Printf("Failed to listen and server %s", server.config.Addr)
	}

	return nil
}

func syslogEndpoint(server *Server) {
	log.Print("Setting up syslog endpoint")
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	syslogServer := syslog.NewServer()
	syslogServer.SetFormat(&RFC5424Formatter{})
	syslogServer.SetHandler(handler)
	syslogServer.ListenTCP(server.config.SyslogAddr)
	syslogServer.Boot()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			if logParts["message"] != "" {
				log.Println(logParts["message"])
			}
		}
	}(channel)

	syslogServer.Wait()
}
