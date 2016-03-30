package server

import (
	"bufio"
	"bytes"
	"github.com/crewjam/rfc5424"
	"github.com/jeromer/syslogparser"
	"strings"
	"time"
)

type RFC5424Formatter struct{}

func (f *RFC5424Formatter) GetParser(line []byte) syslogparser.LogParser {
	return NewParser(line)
}

func (f *RFC5424Formatter) GetSplitFunc() bufio.SplitFunc {
	return nil
}

type Parser struct {
	buff    []byte
	message rfc5424.Message
}

func NewParser(buff []byte) *Parser {
	return &Parser{
		buff:    buff,
		message: rfc5424.Message{},
	}
}

func (p *Parser) Location(location *time.Location) {
	// Ignore as RFC5424 syslog always has a timezone
}

func (p *Parser) Parse() error {
	r := bytes.NewReader(p.buff)
	_, err := p.message.ReadFrom(r)
	return err
}

func (p *Parser) Dump() syslogparser.LogParts {
	message := string(p.message.Message)
	message = strings.Replace(message, "\\n", "\n", -1)

	if strings.HasSuffix(message, "\n") {
		message = strings.TrimSuffix(message, "\n")
	}

	return syslogparser.LogParts{
		"priority":        p.message.Priority,
		"facility":        p.message.Priority / 8,
		"severity":        p.message.Priority % 8,
		"version":         "",
		"timestamp":       p.message.Timestamp,
		"hostname":        p.message.Hostname,
		"app_name":        p.message.AppName,
		"proc_id":         p.message.ProcessID,
		"msg_id":          p.message.MessageID,
		"structured_data": p.message.StructuredData,
		"message":         message,
	}
}
