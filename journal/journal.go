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
	writer JournalWriter
	stream []interface{}
}

type JournalWriter interface {
	Write(envelope *envelope) error
}

func NewJournal(writer JournalWriter) Journal {
	j := &journal{writer: writer}
	return j
}

func (journal journal) Record(record interface{}) error {
	envelope := &envelope{
		Event:     record,
		Time:      time.Now(),
		EventType: reflect.TypeOf(record).Name(),
	}
	return journal.writer.Write(envelope)
}
