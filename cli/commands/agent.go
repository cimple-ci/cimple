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
		},
		Action: func(c *cli.Context) {
			logging.SetDefaultLogger("Agent", os.Stdout)

			agentConfig, err := agent.DefaultConfig()
			agentConfig.ServerAddr = c.String("server-addr")
			agentConfig.ServerPort = c.String("server-port")
			agentConfig.EnableTLS = !c.Bool("no-tls")

			CA_Pool := x509.NewCertPool()
			caFile := c.String("tls-ca-file")
			if len(caFile) != 0 {
				caCert, err := ioutil.ReadFile(caFile)
				if err != nil {
					log.Fatal("Could not load CA certificate!")
				}

				CA_Pool.AppendCertsFromPEM(caCert)
			}

			agentConfig.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: c.Bool("skip-tls-verification"),
				RootCAs:            CA_Pool,
			}
			agentConfig.SyslogUrl = fmt.Sprintf("%s:1514", agentConfig.ServerAddr)

			syslogWriter, err := buildSyslogLogger(agentConfig)
			writers := []io.Writer{os.Stdout}

			if err != nil {
				log.Fatalf("Unable to connect to server syslog %s = %+v", agentConfig.SyslogUrl, err)
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
			}

			err = agent.Start()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}

func buildSyslogLogger(config *agent.Config) (*syslog.Writer, error) {
	log.Printf("Attempting to connect to Cimple Server syslog endpoint - %s", config.SyslogUrl)

	return syslog.Dial("tcp", config.SyslogUrl, syslog.LOG_INFO, "Agent", nil)
}
