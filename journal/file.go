package journal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type fileJournalWriter struct {
	path string
}

func (writer *fileJournalWriter) Write(envelope *envelope) error {
	dir := filepath.Dir(writer.path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(writer.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	a, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	if _, err = f.WriteString(string(a) + "\n"); err != nil {
		return err
	}

	f.Sync()

	return nil
}

func NewFileJournalWriter(path string) *fileJournalWriter {
	return &fileJournalWriter{
		path: path,
	}
}
