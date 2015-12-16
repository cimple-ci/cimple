package build

import (
	"github.com/lukesmith/cimple/journal"
	"github.com/lukesmith/cimple/project"
	"github.com/lukesmith/cimple/vcs"
	"io"
)

type BuildConfig struct {
	BuildId       string
	ExplicitTasks []string
	logWriter     io.Writer
	journal       journal.Journal
	project       project.Project
	tasks         map[string]*project.Task
	repoInfo      vcs.VcsInformation
}

func NewBuildConfig(buildId string, logWriter io.Writer, journal journal.Journal, cfg *project.Config, ri vcs.VcsInformation) *BuildConfig {
	return &BuildConfig{
		BuildId:   buildId,
		logWriter: logWriter,
		journal:   journal,
		project:   cfg.Project,
		tasks:     cfg.Tasks,
		repoInfo:  ri,
	}
}
