package config

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
		},
		Tasks: map[string]Task{
			"echo": Task{
				Description: "Description of the echo task",
				Name:        "echo",
				Archive:     []string{"cow.txt"},
				Commands: map[string]Command{
					"echo_hello_world": Command{
						Command: "echo",
						Args:    []string{"hello world"},
					},
					"echo": Command{
						Command: "echo",
						Args:    []string{"moo >> cow.txt"},
					},
					"cat": Command{
						Command: "cat",
						Args:    []string{"cow.txt"},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("file: %s\n\n%#v\n\n%#v", file, actual, expected)
	}
}
