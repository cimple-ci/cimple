package journal

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
)

type JournalFormatter interface {
	Format(envelope *envelope) (string, error)
}

type journalWriter struct {
	writer    io.Writer
	formatter JournalFormatter
}

func NewJournalWriter(writer io.Writer, formatter JournalFormatter) *journalWriter {
	return &journalWriter{
		writer:    writer,
		formatter: formatter,
	}
}

func (writer *journalWriter) Write(envelope *envelope) error {
	formatted, err := writer.formatter.Format(envelope)
	if err != nil {
		return err
	}
	_, err = writer.writer.Write([]byte(formatted))
	return err
}

type jsonFormatter struct{}

func NewJsonFormatter() *jsonFormatter {
	return &jsonFormatter{}
}

func (f *jsonFormatter) Format(envelope *envelope) (string, error) {
	a, err := json.Marshal(envelope)
	if err != nil {
		return "", err
	}
	return fmt.Sprintln(string(a)), nil
}

type textFormatter struct{}

func NewTextFormatter() *textFormatter {
	return &textFormatter{}
}

func (f *textFormatter) Format(envelope *envelope) (string, error) {
	c := color.New(color.FgBlue).SprintFunc()
	s := fmt.Sprintf("[%s - %s]\n", envelope.Time, envelope.EventType)
	return c(s), nil
}
