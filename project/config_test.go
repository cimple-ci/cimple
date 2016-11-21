package project

import (
	"fmt"
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
						Command: "echo",
						Args:    []string{"hello world"},
						Env:     map[string]string{},
					},
					"echo": Command{
						Command: "echo",
						Args:    []string{"moo >> cow.txt"},
						Skip:    true,
						Env:     map[string]string{},
					},
					"cat": Command{
						Command: "cat",
						Args:    []string{"cow.txt"},
						Env: map[string]string{
							"env": "test",
						},
					},
					"scriptfile": Script{
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
