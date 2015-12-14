package journal

import (
	"reflect"
	"time"
)

type record interface {
}

type Journal interface {
	Record(record record) error
}

type envelope struct {
	Event     record    `json:"event"`
	Time      time.Time `json:"time"`
	EventType string    `json:"type"`
}

type journal struct {
	writer JournalWriter
	stream []record
}

type JournalWriter interface {
	Write(envelope *envelope) error
}

func NewJournal(writer JournalWriter) Journal {
	j := &journal{writer: writer}
	return j
}

func (journal journal) Record(record record) error {
	envelope := &envelope{
		Event:     record,
		Time:      time.Now(),
		EventType: reflect.TypeOf(record).Name(),
	}
	return journal.writer.Write(envelope)
}
