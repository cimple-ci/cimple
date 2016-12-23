package cli

import (
	"log"

	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/lukesmith/cimple/agent"
	"github.com/lukesmith/cimple/logging"
	"github.com/lukesmith/syslog"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"os"
)

func Agent() cli.Command {
	return cli.Command{
		Name:    "agent",
		Aliases: []string{"a"},
		Usage:   "Start the Cimple agent",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "server-addr",
				Usage: "The Cimple server address.",
				Value: "localhost",
			},
			cli.StringFlag{
				Name:  "server-port",
				Usage: "The Cimple server port.",
				Value: "8080",
			},
			cli.BoolFlag{
				Name:  "no-tls",
				Usage: "Disable TLS to connect to the server",
			},
			cli.BoolFlag{
				Name:  "skip-tls-verification",
				Usage: "Skip verification of the server TLS certificate",
			},
			cli.StringFlag{
				Name:  "tls-ca-file",
				Usage: "Specifies the path to the server CA certificate file",
			},
			cli.StringSliceFlag{
				Name:  "tag",
				Usage: "Specify tags for the agent",
			},
		},
		Action: func(c *cli.Context) error {
			logging.SetDefaultLogger("Agent", os.Stdout)

			agentConfig, err := agent.DefaultConfig()
			agentConfig.ServerAddr = c.String("server-addr")
			agentConfig.ServerPort = c.String("server-port")
			agentConfig.EnableTLS = !c.Bool("no-tls")

			if agentConfig.EnableTLS == true {
				caFile := c.String("tls-ca-file")
				tlsConfig, err := createTLSConfig(caFile, c.Bool("skip-tls-verification"))
				if err != nil {
					return err
				}
				agentConfig.TLSClientConfig = tlsConfig
			}

			agentConfig.SyslogUrl = fmt.Sprintf("%s:1514", agentConfig.ServerAddr)

			syslogWriter, err := buildSyslogLogger(agentConfig)
			if err != nil {
				return err
			}
			writers := []io.Writer{os.Stdout}

			if err != nil {
				log.Fatalf("Unable to connect to server syslog %s = %+v", agentConfig.SyslogUrl, err)
				return err
			} else {
				writers = append(writers, syslogWriter)
			}
			defer syslogWriter.Close()

			logWriter := io.MultiWriter(writers...)
			logging.SetDefaultLogger("Agent", logWriter)
			logger := logging.CreateLogger("Agent", logWriter)

			agent, err := agent.NewAgent(agentConfig, logger)
			if err != nil {
				log.Fatal(err)
				return err
			}

			err = agent.Start()
			if err != nil {
				log.Fatal(err)
				return err
			}

			return nil
		},
	}
}

func createTLSConfig(caFile string, skipVerify bool) (*tls.Config, error) {
	CA_Pool := x509.NewCertPool()
	if len(caFile) != 0 {
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatal("Could not load CA certificate!")
			return nil, err
		}

		CA_Pool.AppendCertsFromPEM(caCert)
	}

	return &tls.Config{
		InsecureSkipVerify: skipVerify,
		RootCAs:            CA_Pool,
	}, nil
}

func buildSyslogLogger(config *agent.Config) (*syslog.Writer, error) {
	log.Printf("Attempting to connect to Cimple Server syslog endpoint - %s", config.SyslogUrl)
	if config.EnableTLS == true {
		log.Printf("Connecting to syslog endpoint with TLS enabled")
		return syslog.Dial("tcp", config.SyslogUrl, syslog.LOG_INFO, "Agent", config.TLSClientConfig)
	} else {
		log.Printf("Connecting to syslog endpoint with TLS disabled")
		return syslog.Dial("tcp", config.SyslogUrl, syslog.LOG_INFO, "Agent", nil)
	}
}
