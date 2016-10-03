package journal

import (
	"reflect"
	"time"
)

type Journal interface {
	Record(record interface{}) error
}

type envelope struct {
	Event     interface{} `json:"event"`
	Time      time.Time   `json:"time"`
	EventType string      `json:"type"`
}

type journal struct {
	writers []JournalWriter
	stream  []interface{}
}

type JournalWriter interface {
	Write(envelope *envelope) error
}

func NewJournal(writers []JournalWriter) Journal {
	j := &journal{writers: writers}
	return j
}

func (journal journal) Record(record interface{}) error {
	envelope := &envelope{
		Event:     record,
		Time:      time.Now(),
		EventType: reflect.TypeOf(record).Name(),
	}

	for _, writer := range journal.writers {
		writer.Write(envelope)
	}

	return nil
}
