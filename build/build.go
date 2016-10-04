package build

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	exec "os/exec"
	"text/template"

	"fmt"
	"github.com/lukesmith/cimple/env"
	"github.com/lukesmith/cimple/journal"
	"github.com/lukesmith/cimple/logging"
	"github.com/lukesmith/cimple/project"
	"github.com/lukesmith/cimple/vcs"
	"os"
)

type StepVars struct {
	Cimple     *env.CimpleEnvironment
	Project    project.Project
	Vcs        vcs.VcsInformation
	TaskName   string
	WorkingDir string
	HostEnv    map[string]string
	StepEnv    map[string]string
}

func (sv *StepVars) Map() map[string]string {
	m := make(map[string]string)
	m = merge(m, sv.HostEnv)

	m["CIMPLE_VERSION"] = sv.Cimple.Version
	m["CIMPLE_PROJECT_NAME"] = sv.Project.Name
	m["CIMPLE_PROJECT_VERSION"] = sv.Project.Version
	m["CIMPLE_TASK_NAME"] = sv.TaskName
	m["CIMPLE_WORKING_DIR"] = sv.WorkingDir
	m["CIMPLE_VCS"] = sv.Vcs.Vcs
	m["CIMPLE_VCS_BRANCH"] = sv.Vcs.Branch
	m["CIMPLE_VCS_REVISION"] = sv.Vcs.Revision
	m["CIMPLE_VCS_REMOTE_URL"] = sv.Vcs.RemoteUrl
	m["CIMPLE_VCS_REMOTE_NAME"] = sv.Vcs.RemoteName

	m = merge(m, sv.StepEnv)

	return m
}

type StepContext struct {
	Id     string
	Env    *StepVars
	Cmd    string
	Args   []string
	logger *log.Logger
}

func newStepContext(stepId string, taskEnvs map[string]string, stepConfig project.Step) *StepContext {
	stepContext := new(StepContext)
	stepContext.Id = stepId
	stepContext.Env = new(StepVars)

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

	wd, _ := os.Getwd()
	stepContext.Env.WorkingDir = wd
	stepContext.Env.Cimple = env.Cimple()
	stepContext.Env.StepEnv = merge(taskEnvs, stepConfig.GetEnv())
	stepContext.Env.HostEnv = env.EnvironmentVariables()

	return stepContext
}

func (step *StepContext) templatedArgs() ([]string, error) {
	args := []string{}
	for _, v := range step.Args {
		tmpl, err := template.New("t").Parse(v)
		if err != nil {
			return nil, err
		}
		var doc bytes.Buffer
		tmpl.Execute(&doc, step.Env)
		args = append(args, doc.String())
	}

	return args, nil
}

func (step *StepContext) templatedEnvs() (map[string]string, error) {
	env := make(map[string]string)

	for k, v := range step.Env.Map() {
		tmpl, err := template.New("t").Parse(v)
		if err != nil {
			return nil, err
		}
		var doc bytes.Buffer
		tmpl.Execute(&doc, step.Env)
		env[k] = doc.String()
	}

	return env, nil
}

func (step *StepContext) execute(journal journal.Journal, stdout io.Writer, stderr io.Writer) error {
	args, err := step.templatedArgs()
	if err != nil {
		return err
	}

	var cmd = exec.Command(step.Cmd, args...)

	// Clear out env for command so not to inherit current process's environment.
	cmd.Env = []string{}

	env, err := step.templatedEnvs()
	if err != nil {
		return err
	}

	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	journal.Record(stepStarted{Id: step.Id, Env: env, Step: step.Cmd, Args: args})

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()
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
		stepContext.Env.TaskName = task.Name
		stepContext.Env.Project = config.project
		stepContext.Env.Vcs = config.repoInfo

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
