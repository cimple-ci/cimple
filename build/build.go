package build

import (
	"bytes"
	"io"
	"log"
	exec "os/exec"
	"strings"
	"text/template"

	"fmt"
	"github.com/lukesmith/cimple/env"
	"github.com/lukesmith/cimple/journal"
	"github.com/lukesmith/cimple/logging"
	"github.com/lukesmith/cimple/project"
	"os"
)

type CommandContext struct {
	Id     string
	Env    map[string]string
	Cmd    string
	Args   []string
	logger *log.Logger
}

func newCommandContext(commandId string, taskEnvs map[string]string, cmdConfig *project.Command) *CommandContext {
	commandContext := new(CommandContext)
	commandContext.Id = commandId
	commandContext.logger = logging.CreateLogger("Command", os.Stdout)
	commandContext.Cmd = cmdConfig.Command
	commandContext.Args = cmdConfig.Args
	commandEnvs := merge(taskEnvs, cmdConfig.Env)
	commandContext.Env = merge(commandEnvs, env.Cimple())
	return commandContext
}

func (command *CommandContext) execute(journal journal.Journal, stdout io.Writer, stderr io.Writer) error {
	command.logger.Printf("Executing %s %s", command.Cmd, strings.Join(command.Args, " "))
	var cmd = exec.Command(command.Cmd, command.Args...)

	// Clear out env for command so not to inherit current process's environment.
	cmd.Env = []string{}

	a := make(map[string]string)

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
		a[k] = doc.String()
	}

	journal.Record(commandStarted{Id: command.Id, Env: a, Command: command.Cmd, Args: command.Args})

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()

	if err != nil {
		journal.Record(commandFailed{Id: command.Id})
		return err
	}

	journal.Record(commandSuccessful{Id: command.Id})

	return nil
}

type BuildTask struct {
	Name     string
	Commands []CommandContext
}

type Build struct {
	ID     int
	Tasks  []BuildTask
	config *BuildConfig
	logger *log.Logger
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func NewBuild(config *BuildConfig) (*Build, error) {
	build := new(Build)
	build.config = config
	build.logger = logging.CreateLogger("Build", config.logWriter)
	build.ID = 1

	for _, task := range config.tasks {
		if len(config.ExplicitTasks) != 0 {
			if contains(config.ExplicitTasks, task.Name) {
				if task.Skip {
					build.logger.Printf("Unskipping task %s as explicitly specified", task.Name)
					task.Skip = false
				}
			} else {
				build.logger.Printf("Skipping task %s as not explicitly specified", task.Name)
				task.Skip = true
			}
		}

		commandContexts, err := buildCommandContexts(build.logger, config, task)
		if err != nil {
			return nil, err
		}
		buildTask := BuildTask{task.Name, commandContexts}
		build.Tasks = append(build.Tasks, buildTask)
	}

	return build, nil
}

func (build *Build) Run() error {
	build.logger.Printf("Running build #%d", build.ID)
	build.config.journal.Record(buildStarted{Repo: build.config.repoInfo})

	for _, task := range build.Tasks {
		commandIds := []string{}

		for _, command := range task.Commands {
			commandIds = append(commandIds, command.Id)
		}

		build.config.journal.Record(taskStarted{Id: task.Name, Commands: commandIds})

		for _, command := range task.Commands {
			err := command.execute(build.config.journal, build.config.logWriter, build.config.logWriter)
			if err != nil {
				build.config.journal.Record(taskFailed{Id: task.Name})
				return err
			}
		}

		build.config.journal.Record(taskSuccessful{Id: task.Name})
	}

	build.config.journal.Record("Build finished successfully")

	return nil
}

func buildCommandContexts(logger *log.Logger, config *BuildConfig, task *project.Task) ([]CommandContext, error) {
	var contexts []CommandContext

	if task.Skip {
		logger.Printf("Skipping task: %s", task.Name)
		return []CommandContext{}, nil
	}

	taskEnvs := merge(config.project.Env, task.Env)

	for k, command := range task.Commands {
		commandId := fmt.Sprintf("%s-%s", task.Name, k)

		if command.Skip {
			config.journal.Record(skipCommand{Id: commandId})
			continue
		}

		commandContext := newCommandContext(commandId, taskEnvs, &command)
		commandContext.logger = logging.CreateLogger("Command", config.logWriter)
		commandContext.Env["CIMPLE_PROJECT_NAME"] = config.project.Name
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
