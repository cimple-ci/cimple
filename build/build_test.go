package build

import (
	"log"
	"os"
	"testing"

	"github.com/lukesmith/cimple/env"
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

func Test_StepVars_Map(t *testing.T) {
	vars := new(StepVars)
	vars.HostEnv = make(map[string]string)
	vars.StepEnv = make(map[string]string)
	vars.HostEnv["HOST_ENV"] = "4212"
	vars.StepEnv["STEP_ENV"] = "1234"
	vars.WorkingDir = "/c/temp"
	vars.TaskName = "taskname"
	vars.Cimple = &env.CimpleEnvironment{
		Version: "1.5.3",
	}
	p := &project.Project{
		Name:    "projectname",
		Version: "4.3.1",
	}
	vars.Project = *p

	m := vars.Map()

	if m["CIMPLE_VERSION"] != "1.5.3" {
		t.Fatalf("Expected CIMPLE_VERSION to be 1.5.3 - was %s", m["CIMPLE_VERSION"])
	}

	if m["CIMPLE_PROJECT_NAME"] != "projectname" {
		t.Fatalf("Expected CIMPLE_PROJECT_NAME to be projeoctname - was %s", m["CIMPLE_PROJECT_NAME"])
	}

	if m["CIMPLE_PROJECT_VERSION"] != "4.3.1" {
		t.Fatalf("Expected CIMPLE_PROJECT_VERSION to be 4.3.1 - was %s", m["CIMPLE_PROJECT_VERSION"])
	}

	if m["CIMPLE_TASK_NAME"] != "taskname" {
		t.Fatalf("Expected CIMPLE_TASK_NAME to be taskname - was %s", m["CIMPLE_TASK_NAME"])
	}

	if m["CIMPLE_WORKING_DIR"] != "/c/temp" {
		t.Fatalf("Expected CIMPLE_WORKING_DIR to be /c/temp - was %s", m["CIMPLE_WORKING_DIR"])
	}

	if m["HOST_ENV"] != "4212" {
		t.Fatalf("Expected HOST_ENV to be 4212 - was %s", m["HOST_ENV"])
	}

	if m["STEP_ENV"] != "1234" {
		t.Fatalf("Expected STEP_ENV to be 1234 - was %s", m["STEP_ENV"])
	}
}

func Test_StepVars_Map_Precedence(t *testing.T) {
	vars := new(StepVars)
	vars.HostEnv = make(map[string]string)
	vars.StepEnv = make(map[string]string)
	vars.HostEnv["CIMPLE_VERSION"] = "1"
	vars.HostEnv["CIMPLE_PROJECT_NAME"] = "a"

	vars.Cimple = &env.CimpleEnvironment{
		Version: "2",
	}
	p := &project.Project{
		Name: "b",
	}
	vars.Project = *p

	vars.StepEnv["CIMPLE_PROJECT_NAME"] = "c"

	m := vars.Map()

	if m["CIMPLE_VERSION"] != "2" {
		t.Fatalf("Expected CIMPLE_VERSION to be overriden from HostEnv")
	}

	if m["CIMPLE_PROJECT_NAME"] != "c" {
		t.Fatalf("Expected CIMPLE_PROJECT_NAME to be overriden from StepEnv")
	}
}

type fakeJournal struct {
}

func (f fakeJournal) Record(record interface{}) error {
	return nil
}
