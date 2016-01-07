package project

import (
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
			Env: map[string]string{
				"project_env": "project",
			},
		},
		Tasks: map[string]*Task{
			"echo": &Task{
				Description: "Description of the echo task",
				Name:        "echo",
				Skip:        true,
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
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("file: %s\n\n%#v\n\n%#v", file, actual, expected)
	}
}

// TODO: Test task names are unique
// TODO: Test step names within a task are unique
