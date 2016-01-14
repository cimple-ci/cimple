package cli

import (
	"flag"
	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/build"
	"github.com/lukesmith/cimple/project"
	"github.com/lukesmith/cimple/vcs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun_Settings(t *testing.T) {
	assert := assert.New(t)

	command := Run()

	assert.Equal("run", command.Name)
	assert.Equal([]string{"r"}, command.Aliases)
}

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

	command := Run()

	set := flag.NewFlagSet("test", flag.ContinueOnError)
	for _, f := range command.Flags {
		f.Apply(set)
	}
	set.Parse([]string{"--task", "two"})
	context := cli.NewContext(nil, set, nil)
	command.Action(context)

	assert.Equal([]string{"two"}, executedConfig.ExplicitTasks)
}