package project

import (
	"os"
	"strings"
)

type ProjectVariables struct {
	WorkingDir string
	HostEnvVar map[string]string
}

func GetVariables() (*ProjectVariables, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pv := new(ProjectVariables)
	pv.WorkingDir = wd
	pv.HostEnvVar = make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		pv.HostEnvVar[pair[0]] = pair[1]
	}

	return pv, nil
}
