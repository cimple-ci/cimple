package project

import (
	"fmt"
	"github.com/lukesmith/cimple/env"
	"github.com/lukesmith/cimple/vcs"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	file := "basic.hcl"
	path, err := filepath.Abs(filepath.Join("./test-fixtures", file))
	if err != nil {
		t.Fatalf("File: %s\n\n%s", file, err)
	}

	actual, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("file: %s\n\n%s", file, err)
	}

	expected := &Config{
		Project: Project{
			Name:        "Cimple",
			Description: "Project description",
			Version:     "1.3.2",
			Env: map[string]string{
				"project_env": "project",
			},
		},
		Tasks: map[string]*Task{
			"echo": &Task{
				Description: "Description of the echo task",
				Depends:     []string{},
				Name:        "echo",
				Skip:        true,
				LimitTo:     "",
				Archive:     []string{"cow.txt"},
				Env: map[string]string{
					"task_env": "global",
				},
				StepOrder: []string{"echo_hello_world", "echo", "scriptfile", "cat"},
				Steps: map[string]Step{
					"echo_hello_world": Command{
						name:    "echo_hello_world",
						Command: "echo",
						Args:    []string{"hello world"},
						Env:     map[string]string{},
					},
					"echo": Command{
						name:    "echo",
						Command: "echo",
						Args:    []string{"moo >> cow.txt"},
						Skip:    true,
						Env:     map[string]string{},
					},
					"cat": Command{
						name:    "cat",
						Command: "cat",
						Args:    []string{"cow.txt"},
						Env: map[string]string{
							"env": "test",
						},
					},
					"scriptfile": Script{
						name: "scriptfile",
						Body: "echo 1",
						Env: map[string]string{
							"env": "test",
						},
					},
				},
			},
			"publish": &Task{
				Description: "Publish packages",
				Depends:     []string{"echo"},
				Name:        "publish",
				Skip:        false,
				Archive:     []string{},
				Env:         map[string]string{},
				StepOrder:   []string{},
				Steps:       map[string]Step{},
				LimitTo:     "server",
			},
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("file: %s\n\n%#v\n\n%#v", file, actual, expected)
	}
}

func TestTaskNamesUnique(t *testing.T) {
	file := "duplicate-task-names.hcl"
	path, err := filepath.Abs(filepath.Join("./test-fixtures", file))
	if err != nil {
		t.Fatalf("File: %s\n\n%s", file, err)
	}

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatalf("file: %s\n\nExpected error due to duplicate task names", file)
	}
}

func TestStepNamesUniqueWithinATask(t *testing.T) {
	file := "duplicate-step-names.hcl"
	path, err := filepath.Abs(filepath.Join("./test-fixtures", file))
	if err != nil {
		t.Fatalf("File: %s\n\n%s", file, err)
	}

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatalf("%s\n\nExpected error due to duplicate step names within a task", file)
	}
}

func TestStepNamesContainValidCharacters(t *testing.T) {
	file := "invalid-step-names.hcl"
	path, err := filepath.Abs(filepath.Join("./test-fixtures", file))
	if err != nil {
		t.Fatalf("File: %s\n\n%s", file, err)
	}

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatalf("file: %s\n\nExpected error due to invalid step names", file)
	}
}

func TestTaskNamesContainValidCharacters(t *testing.T) {
	file := "invalid-task-names.hcl"
	path, err := filepath.Abs(filepath.Join("./test-fixtures", file))
	if err != nil {
		t.Fatalf("File: %s\n\n%s", file, err)
	}

	_, err = LoadConfig(path)
	if err == nil {
		t.Fatalf("file: %s\n\nExpected error due to invalid task names", file)
	}
}

func TestMissingName(t *testing.T) {
	const testconfig = `
	version = "0.0.1"
	`

	_, err := Load(testconfig)
	assert.Equal(t, fmt.Errorf("'name' was not specified"), err)
}

func TestMissingVersion(t *testing.T) {
	const testconfig = `
	name = "test"
	`

	_, err := Load(testconfig)
	assert.Equal(t, fmt.Errorf("'version' was not specified"), err)
}

func TestDescriptionOptional(t *testing.T) {
	const testconfig = `
	name = "test"
	version = "0.0.1"
	`

	_, err := Load(testconfig)
	assert.Nil(t, err)
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
	p := &Project{
		Name:    "projectname",
		Version: "4.3.1",
	}
	vars.Project = *p
	v := &vcs.VcsInformation{
		Vcs:        "git",
		Branch:     "my-branch",
		Revision:   "12345",
		RemoteUrl:  "git@github.com/cimple",
		RemoteName: "origin",
	}
	vars.Vcs = *v

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

	if m["CIMPLE_VCS"] != "git" {
		t.Fatalf("Expected CIMPLE_VCS to be git - was %s", m["CIMPLE_VCS"])
	}

	if m["CIMPLE_VCS_BRANCH"] != "my-branch" {
		t.Fatalf("Expected CIMPLE_VCS_BRANCH to be my-branch - was %s", m["CIMPLE_VCS_BRANCH"])
	}

	if m["CIMPLE_VCS_REVISION"] != "12345" {
		t.Fatalf("Expected CIMPLE_VCS_REVISION to be 12345 - was %s", m["CIMPLE_VCS_REVISION"])
	}

	if m["CIMPLE_VCS_REMOTE_URL"] != "git@github.com/cimple" {
		t.Fatalf("Expected CIMPLE_VCS_REMOTE_URL to be git@github.com/cimple - was %s", m["CIMPLE_VCS_REMOTE_URL"])
	}

	if m["CIMPLE_VCS_REMOTE_NAME"] != "origin" {
		t.Fatalf("Expected CIMPLE_VCS_REMOTE_NAME to be origin - was %s", m["CIMPLE_VCS_REMOTE_NAME"])
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
	p := &Project{
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
