package build

import "io"
import "github.com/lukesmith/cimple/journal"
import (
	"fmt"
	"github.com/lukesmith/cimple/project"
	"log"
)

type buildConfig struct {
	logWriter io.Writer
	journal   journal.Journal
	project   *project.Project
	tasks     map[string]project.Task
}

func NewBuildConfig(logWriter io.Writer, journal journal.Journal, project *project.Project, tasks map[string]project.Task) *buildConfig {
	return &buildConfig{
		logWriter: logWriter,
		journal:   journal,
		project:   project,
		tasks:     tasks,
	}
}

func (cfg *buildConfig) createLogger(prefix string) *log.Logger {
	return log.New(cfg.logWriter, fmt.Sprintf("%s: ", prefix), log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
}
