package logging

import (
	"fmt"
	"io"
	"log"
)

func CreateLogger(prefix string, out io.Writer) *log.Logger {
	return log.New(out, fmt.Sprintf("%-10s: ", prefix), log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
}
