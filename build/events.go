package build

import (
	"github.com/lukesmith/cimple/vcs"
)

type stepStarted struct {
	Id   string
	Env  map[string]string
	Step string
	Args []string
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
