package cli

import (
	"flag"
	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/project"
	"github.com/stretchr/testify/assert"
	"testing"
	"io"
)

func TestRun_Settings(t *testing.T) {
	assert := assert.New(t)

	command := Run()

	assert.Equal("run", command.Name)
	assert.Equal([]string{"r"}, command.Aliases)
}

func TestRun_ExplicitTasks(t *testing.T) {
	assert := assert.New(t)
	var executedConfig *project.Config

	loadConfig = func() (*project.Config, error) {
		return &project.Config{
			Tasks: map[string]*project.Task{
				"one": &project.Task{Name: "one", Skip: false},
				"two": &project.Task{Name: "two", Skip: true},
			},
		}, nil
	}

	executeBuild = func(runId string, out io.Writer, c *cli.Context, cfg *project.Config) error {
		executedConfig = cfg
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

	assert.True(executedConfig.Tasks["one"].Skip)
	assert.False(executedConfig.Tasks["two"].Skip)
}
