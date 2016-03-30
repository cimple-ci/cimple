package journal

import (
	"encoding/json"
	"github.com/lukesmith/syslog"
)

type syslogWriter struct {
	syslog *syslog.Writer
}

func (writer *syslogWriter) Write(envelope *envelope) error {
	a, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	writer.syslog.Info(string(a))

	return nil
}

func NewSyslogWriter(syslog *syslog.Writer) *syslogWriter {
	return &syslogWriter{
		syslog: syslog,
	}
}
