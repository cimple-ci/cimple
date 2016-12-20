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
	Addr        string
	SyslogAddr  string
	TLSCertFile string
	TLSKeyFile  string
	EnableTLS   bool
}

func DefaultConfig() *Config {
	return &Config{
		Addr:        ":0",
		SyslogAddr:  ":1514",
		EnableTLS:   false,
		TLSCertFile: "server.crt",
		TLSKeyFile:  "server.key",
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
		logger.Printf("CertFile %s", config.TLSCertFile)
		logger.Printf("KeyFile %s", config.TLSKeyFile)
		cer, err := tls.LoadX509KeyPair(config.TLSCertFile, config.TLSKeyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cer},
		}

		return tls.Listen("tcp", addr, tlsConfig)
	} else {
		logger.Printf("Configuring server with TLS disabled")
		return net.Listen("tcp", addr)
	}
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
