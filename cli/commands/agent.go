package cli

import (
	"log"

	"fmt"
	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/agent"
	"github.com/lukesmith/cimple/logging"
	"github.com/lukesmith/syslog"
	"io"
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
		},
		Action: func(c *cli.Context) {
			agentConfig, err := agent.DefaultConfig()
			agentConfig.ServerAddr = c.String("server-addr")
			agentConfig.ServerPort = c.String("server-port")
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
