package runner

import (
	"github.com/lukesmith/cimple/build"
	"github.com/lukesmith/cimple/project"
	"github.com/lukesmith/cimple/vcs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun_ExplicitTasks(t *testing.T) {
	assert := assert.New(t)
	var executedConfig *build.BuildConfig

	loadConfig = func() (*project.Config, error) {
		return &project.Config{
			Tasks: map[string]*project.Task{
				"one": &project.Task{Name: "one", Skip: false},
				"two": &project.Task{Name: "two", Skip: true},
			},
		}, nil
	}

	loadRepositoryInfo = func() *vcs.VcsInformation {
		return new(vcs.VcsInformation)
	}

	executeBuild = func(buildConfig *build.BuildConfig) error {
		executedConfig = buildConfig
		return nil
	}

	options := &RunOptions{
		Journal: &JournalSettings{},
	}
	Run(options, []string{"two"})

	assert.Equal([]string{"two"}, executedConfig.ExplicitTasks)
}

func TestRun_RunContext(t *testing.T) {
	assert := assert.New(t)
	var executedConfig *build.BuildConfig

	loadConfig = func() (*project.Config, error) {
		return &project.Config{
			Tasks: map[string]*project.Task{
				"one": &project.Task{Name: "one", Skip: false},
			},
		}, nil
	}

	loadRepositoryInfo = func() *vcs.VcsInformation {
		return new(vcs.VcsInformation)
	}

	executeBuild = func(buildConfig *build.BuildConfig) error {
		executedConfig = buildConfig
		return nil
	}

	options := &RunOptions{
		Context: "xyz",
		Journal: &JournalSettings{},
	}
	Run(options, []string{})

	assert.Equal(options.Context, executedConfig.RunContext)
}
