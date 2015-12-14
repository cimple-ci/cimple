package journal

import "testing"

func TestJournalFileRecord(t *testing.T) {
	writer := &testWriter{}
	journal := NewJournal(writer)
	event := &testEvent{}
	err := journal.Record(event)
	if err != nil {
		t.Fatalf("Failed to record event %s", event, err)
	}

	if len(writer.written) != 1 {
		t.Fatalf("Expected a single event to have been written")
	}

	if writer.written[0].Event != event {
		t.Fatalf("Exepected %s to have been written", event, err)
	}
}

type testWriter struct {
	out     JournalWriter
	written []envelope
}

func (writer *testWriter) Write(envelope *envelope) error {
	writer.written = append(writer.written, *envelope)
	return nil
}

type testEvent struct {
}
