package build

import (
	"io"
	"github.com/lukesmith/cimple/journal"
	"github.com/lukesmith/cimple/project"
)

type buildConfig struct {
	logWriter io.Writer
	journal   journal.Journal
	project   *project.Project
	tasks     map[string]*project.Task
}

func NewBuildConfig(logWriter io.Writer, journal journal.Journal, project *project.Project, tasks map[string]*project.Task) *buildConfig {
	return &buildConfig{
		logWriter: logWriter,
		journal:   journal,
		project:   project,
		tasks:     tasks,
	}
}
