package logging

import (
	"fmt"
	"io"
	"log"
)

func SetDefaultLogger(prefix string, out io.Writer) {
	log.SetOutput(out)
	log.SetPrefix(fmt.Sprintf("%-10s: ", prefix))
	log.SetFlags(log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
}

func CreateLogger(prefix string, out io.Writer) *log.Logger {
	return log.New(out, fmt.Sprintf("%-10s: ", prefix), log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
}
