package build

import (
	"log"
	"os"
	"testing"

	"github.com/lukesmith/cimple/project"
	"github.com/lukesmith/cimple/vcs"
)

func Test_buildStepContexts_OrdersContexts(t *testing.T) {
	var task = project.Task{
		Name: "bob",
	}
	task.StepOrder = []string{"firststep", "secondstep"}
	task.Steps = map[string]project.Step{
		"secondstep": project.Command{
			Command: "echo",
			Args:    []string{"hello world"},
			Env:     map[string]string{},
		},
		"firststep": project.Command{
			Command: "echo",
			Args:    []string{"moo >> cow.txt"},
			Env:     map[string]string{},
		},
	}

	project := project.Config{}
	journal := fakeJournal{}
	vcs := vcs.VcsInformation{}
	var buildConfig = NewBuildConfig("test", os.Stdout, journal, &project, vcs)
	var logger = log.New(os.Stdout, "test", log.LUTC)

	contexts, err := buildStepContexts(logger, buildConfig, &task)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if contexts[0].Id != "bob.firststep" {
		t.Fatalf("firststep is not the first context - %s", contexts[0].Id)
	}

	if contexts[1].Id != "bob.secondstep" {
		t.Fatalf("secondstep is not the second context - %s", contexts[1].Id)
	}
}

func Test_buildStepContexts_SkipTask(t *testing.T) {
	var task = project.Task{}
	task.Skip = true
	var buildConfig = BuildConfig{}
	var logger = log.New(os.Stdout, "test", log.LUTC)

	contexts, err := buildStepContexts(logger, &buildConfig, &task)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(contexts) != 0 {
		t.Fatalf("Task was not skipped")
	}
}

type fakeJournal struct {
}

func (f fakeJournal) Record(record interface{}) error {
	return nil
}
