package project

import (
	"os"
)

type ProjectVariables struct {
	WorkingDir string
}

func GetProjectVariables() (*ProjectVariables, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	e := new(ProjectVariables)
	e.WorkingDir = wd

	return e, nil
}
