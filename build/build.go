package build

import (
	"bytes"
	"github.com/lukesmith/cimple/env"
	"github.com/lukesmith/cimple/project"
	"io"
	"log"
	exec "os/exec"
	"strings"
	"text/template"
)

type CommandContext struct {
	Env  map[string]string
	Cmd  string
	Args []string
}

func (command *CommandContext) execute(stdout io.Writer, stderr io.Writer) error {
	log.Printf("Executing %s %s", command.Cmd, strings.Join(command.Args, " "))
	var cmd = exec.Command(command.Cmd, command.Args...)

	// Clear out env for command so not to inherit current process's environment.
	cmd.Env = []string{}

	for k, v := range command.Env {
		// TODO: Extract templating environment variables
		tmpl, err := template.New("t").Parse(v)
		if err != nil {
			return err
		}
		vari, err := project.GetVariables()
		var doc bytes.Buffer
		tmpl.Execute(&doc, vari)
		cmd.Env = append(cmd.Env, k+"="+doc.String())
	}

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Run()

	return nil
}

type BuildTask struct {
	Commands []CommandContext
}

type Build struct {
	ID    int
	Tasks []BuildTask
}

func NewBuild(config *project.Config) (*Build, error) {
	build := new(Build)
	build.ID = 1

	for _, task := range config.Tasks {
		commandContexts, err := buildCommandContexts(&config.Project, &task)
		if err != nil {
			return nil, err
		}
		buildTask := BuildTask{commandContexts}
		build.Tasks = append(build.Tasks, buildTask)
	}

	return build, nil
}

func (build *Build) Run(stdout io.Writer, stderr io.Writer) error {
	log.Printf("Running build #%d", build.ID)

	for _, task := range build.Tasks {
		for _, command := range task.Commands {
			err := command.execute(stdout, stderr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func buildCommandContexts(project *project.Project, task *project.Task) ([]CommandContext, error) {
	var contexts []CommandContext

	if task.Skip {
		log.Printf("Skipping task: %s", task.Name)
		return []CommandContext{}, nil
	}

	taskEnvs := merge(project.Env, task.Env)

	for k, command := range task.Commands {
		if command.Skip {
			log.Printf("Skipping task: %s command: %s", task.Name, k)
			continue
		}

		commandContext := new(CommandContext)
		commandContext.Cmd = command.Command
		commandContext.Args = command.Args
		commandEnvs := merge(taskEnvs, command.Env)
		commandContext.Env = merge(commandEnvs, env.Cimple())
		commandContext.Env["CIMPLE_PROJECT_NAME"] = project.Name
		commandContext.Env["CIMPLE_TASK_NAME"] = task.Name

		contexts = append(contexts, *commandContext)
	}

	return contexts, nil
}

func merge(a map[string]string, b map[string]string) map[string]string {
	c := make(map[string]string)

	for k, v := range a {
		c[k] = v
	}

	for k, v := range b {
		c[k] = v
	}

	return c
}
