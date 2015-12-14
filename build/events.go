package build

type commandStarted struct {
	Id      string
	Env     map[string]string
	Command string
	Args    []string
}

type commandSuccessful struct {
	Id string
}

type commandFailed struct {
	Id string
}

type skipCommand struct {
	Id string
}

type taskStarted struct {
	Id       string
	Commands []string
}

type taskFailed struct {
	Id string
}

type taskSuccessful struct {
	Id string
}
