package server

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/lukesmith/cimple/database"
	"github.com/mcuadros/go-syslog"
)

type Config struct {
	Addr            string
	SyslogAddr      string
	EnableTLS       bool
	TLSServerConfig *tls.Config
}

func DefaultConfig() *Config {
	return &Config{
		Addr:            ":0",
		SyslogAddr:      ":1514",
		EnableTLS:       false,
		TLSServerConfig: &tls.Config{},
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

	ln, err := createListener(server.config, server.logger)
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

func createListener(config *Config, logger *log.Logger) (net.Listener, error) {
	addr := config.Addr
	if addr == "" {
		addr = ":http"
	}

	if config.EnableTLS {
		logger.Printf("Configuring server with TLS enabled")
		return tls.Listen("tcp", addr, config.TLSServerConfig)
	} else {
		logger.Printf("Configuring server with TLS disabled")
		return net.Listen("tcp", addr)
	}
}

func syslogEndpoint(server *Server) error {
	server.logger.Printf("Setting up syslog endpoint at %s", server.config.SyslogAddr)
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	syslogServer := syslog.NewServer()
	syslogServer.SetFormat(&RFC5424Formatter{})
	syslogServer.SetHandler(handler)
	syslogServer.SetTlsPeerNameFunc(func(tlsConn *tls.Conn) (string, bool) {
		return "", true
	})
	if server.config.EnableTLS == true {
		server.logger.Print("Syslog listening with TLS enabled")
		err := syslogServer.ListenTCPTLS(server.config.SyslogAddr, server.config.TLSServerConfig)
		if err != nil {
			return err
		}
	} else {
		server.logger.Print("Syslog listening with TLS disabled")
		err := syslogServer.ListenTCP(server.config.SyslogAddr)
		if err != nil {
			return err
		}
	}

	err := syslogServer.Boot()
	if err != nil {
		return err
	}

	server.logger.Printf("Syslog server booted")

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			if logParts["message"] != "" {
				server.logger.Println(logParts["message"])
			}
		}
	}(channel)

	syslogServer.Wait()

	return nil
}
