package server

import (
	"log"
	"net"
	"net/http"

	"github.com/lukesmith/cimple/database"
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

func NewServer(config *Config, logger *log.Logger) (*Server, error) {
	s := &Server{
		config: config,
		logger: logger,
	}

	return s, nil
}

func (server *Server) Start() error {
	agentPool := newAgentPool(server.logger)
	bq := &buildQueue{}
	bq.queue = make(chan interface{})
	bq.agentpool = agentPool

	db := database.NewDatabase("./.cimple")
	app := NewFrontend(db, agentPool, bq, server.config.Addr, server.logger)

	http.Handle("/", app)

	go syslogEndpoint(server)

	go agentPool.run()
	go bq.run()

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
	server.logger.Print("Setting up syslog endpoint")
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
				server.logger.Println(logParts["message"])
			}
		}
	}(channel)

	syslogServer.Wait()
}
