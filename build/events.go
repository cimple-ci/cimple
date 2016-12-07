package build

import (
	"github.com/lukesmith/cimple/project"
	"github.com/lukesmith/cimple/vcs"
)

type stepStarted struct {
	Id       string
	Env      *project.StepVars //map[string]string
	StepType string
	Step     interface{}
}

type stepSuccessful struct {
	Id string
}

type stepFailed struct {
	Id string
}

type skipStep struct {
	Id string
}

type taskStarted struct {
	Id    string
	Steps []string
}

type taskSkipped struct {
	Id     string
	Reason string
}

type taskFailed struct {
	Id string
}

type taskSuccessful struct {
	Id string
}

type buildStarted struct {
	Repo vcs.VcsInformation
}
