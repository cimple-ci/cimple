package build

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
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

type StepContext struct {
	Id     string
	Env    map[string]string
	Cmd    string
	Args   []string
	logger *log.Logger
}

func newStepContext(stepId string, taskEnvs map[string]string, stepConfig project.Step) *StepContext {
	stepContext := new(StepContext)
	stepContext.Id = stepId

	command, ok := stepConfig.(project.Command)
	if ok {
		stepContext.Cmd = command.Command
		stepContext.Args = command.Args
	}

	script, ok := stepConfig.(project.Script)
	if ok {
		f, _ := ioutil.TempFile(os.TempDir(), "step")
		f.WriteString(script.Body)
		f.Close()

		stepContext.Cmd = "/bin/sh"
		stepContext.Args = []string{f.Name()}
	}

	stepEnvs := merge(taskEnvs, stepConfig.GetEnv())
	stepContext.Env = merge(stepEnvs, env.Cimple())
	return stepContext
}

func (step *StepContext) execute(journal journal.Journal, stdout io.Writer, stderr io.Writer) error {
	step.logger.Printf("Executing %s %s", step.Cmd, strings.Join(step.Args, " "))
	var cmd = exec.Command(step.Cmd, step.Args...)

	// Clear out env for command so not to inherit current process's environment.
	cmd.Env = []string{}

	a := make(map[string]string)

	for k, v := range step.Env {
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

	journal.Record(stepStarted{Id: step.Id, Env: a, Step: step.Cmd, Args: step.Args})

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()

	if err != nil {
		journal.Record(stepFailed{Id: step.Id})
		return err
	}

	journal.Record(stepSuccessful{Id: step.Id})

	return nil
}

type BuildTask struct {
	Name  string
	Steps []StepContext
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

		contexts, err := buildStepContexts(build.logger, config, task)
		if err != nil {
			return nil, err
		}
		buildTask := BuildTask{task.Name, contexts}
		build.Tasks = append(build.Tasks, buildTask)
	}

	return build, nil
}

func (build *Build) Run() error {
	build.logger.Printf("Running build #%d", build.ID)
	build.config.journal.Record(buildStarted{Repo: build.config.repoInfo})

	for _, task := range build.Tasks {
		stepIds := []string{}

		for _, step := range task.Steps {
			stepIds = append(stepIds, step.Id)
		}

		build.config.journal.Record(taskStarted{Id: task.Name, Steps: stepIds})

		for _, step := range task.Steps {
			err := step.execute(build.config.journal, build.config.logWriter, build.config.logWriter)
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

func buildStepContexts(logger *log.Logger, config *BuildConfig, task *project.Task) ([]StepContext, error) {
	var contexts []StepContext

	if task.Skip {
		logger.Printf("Skipping task: %s", task.Name)
		return []StepContext{}, nil
	}

	taskEnvs := merge(config.project.Env, task.Env)

	for _, stepName := range task.StepOrder {
		step, found := task.Steps[stepName]

		if !found {
			// TODO: include list of possible step names in error
			return []StepContext{}, errors.New(fmt.Sprintf("Could not find step named %s.", stepName))
		}

		stepId := fmt.Sprintf("%s.%s", task.Name, stepName)

		if step.GetSkip() {
			config.journal.Record(skipStep{Id: stepId})
			continue
		}

		stepContext := newStepContext(stepId, taskEnvs, step)
		stepContext.logger = logging.CreateLogger("Step", config.logWriter)
		stepContext.Env["CIMPLE_PROJECT_NAME"] = config.project.Name
		stepContext.Env["CIMPLE_TASK_NAME"] = task.Name

		contexts = append(contexts, *stepContext)
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
