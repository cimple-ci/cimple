package build

import (
	"errors"
	"log"
	"reflect"
	"time"

	"fmt"
	"github.com/lukesmith/cimple/env"
	"github.com/lukesmith/cimple/logging"
	"github.com/lukesmith/cimple/project"
	"os"
)

type StepContext struct {
	Id     string
	Env    *project.StepVars
	Cmd    string
	logger *log.Logger
	Step   project.Step
}

func newStepContext(stepId string, taskEnvs map[string]string, stepConfig project.Step) *StepContext {
	stepContext := new(StepContext)
	stepContext.Id = stepId
	stepContext.Env = new(project.StepVars)

	wd, _ := os.Getwd()
	stepContext.Env.BuildDate = time.Now()
	stepContext.Env.WorkingDir = wd
	stepContext.Env.Cimple = env.Cimple()
	stepContext.Env.StepEnv = merge(taskEnvs, stepConfig.GetEnv())
	stepContext.Env.HostEnv = env.EnvironmentVariables()
	stepContext.Step = stepConfig

	return stepContext
}

type BuildTask struct {
	Name         string
	Steps        []StepContext
	dependencies []string
	limitTo      string
	skip         bool
}

func (bt BuildTask) GetID() string {
	return bt.Name
}

func (bt BuildTask) GetDependencies() []string {
	return bt.dependencies
}

type Build struct {
	ID     int
	tasks  map[string]*BuildTask
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
	build.tasks = make(map[string]*BuildTask)

	for _, task := range config.tasks {
		contexts, err := buildStepContexts(build.logger, build.config, task)
		if err != nil {
			return nil, err
		}

		buildTask := &BuildTask{
			Name:         task.Name,
			Steps:        contexts,
			skip:         task.Skip,
			dependencies: task.Depends,
			limitTo:      task.LimitTo,
		}
		build.tasks[task.Name] = buildTask
	}

	return build, nil
}

func (build *Build) checkSkip(task *BuildTask) (string, bool) {
	if len(build.config.ExplicitTasks) != 0 {
		if contains(build.config.ExplicitTasks, task.Name) {
			if task.skip {
				build.logger.Printf("Unskipping task %s as explicitly specified", task.Name)
				return "", false
			}
		} else {
			build.logger.Printf("Skipping task %s as not explicitly specified", task.Name)
			return "Explicit tasks defined", true
		}
	}

	if len(task.limitTo) != 0 && task.limitTo != build.config.RunContext {
		build.logger.Printf("Skipping task %s. Is limited to run in %s context. Current context is %s", task.Name, task.limitTo, build.config.RunContext)
		return "Outside of run context", true
	}

	return "", false
}

func (build *Build) Run() error {
	build.logger.Printf("Running build #%d", build.ID)
	build.config.journal.Record(buildStarted{Repo: build.config.repoInfo})

	tasks := []TaskNode{}
	for _, t := range build.tasks {
		tasks = append(tasks, t)
	}

	buildStrategy := NewBuildStrategy(tasks)
	err := buildStrategy.Build(func(taskName string) error {
		return build.runTask(build.tasks[taskName])
	})
	if err != nil {
		return err
	}

	build.config.journal.Record("Build finished successfully")

	return nil
}

func (build *Build) runTask(task *BuildTask) error {
	if reason, skip := build.checkSkip(task); skip {
		build.config.journal.Record(taskSkipped{Id: task.Name, Reason: reason})
		return nil
	}

	build.logger.Printf("Running task %s", task.Name)
	stepIds := []string{}

	for _, step := range task.Steps {
		stepIds = append(stepIds, step.Id)
	}

	build.config.journal.Record(taskStarted{Id: task.Name, Steps: stepIds})

	for _, stepContext := range task.Steps {
		stepType := reflect.TypeOf(stepContext.Step).Name()
		build.config.journal.Record(stepStarted{Id: stepContext.Id, Env: stepContext.Env, StepType: stepType, Step: stepContext.Step})
		err := stepContext.Step.Execute(*stepContext.Env, build.config.logWriter, build.config.logWriter)
		if err != nil {
			build.config.journal.Record(stepFailed{Id: stepContext.Id})
			build.config.journal.Record(taskFailed{Id: task.Name})
			return err
		}

		build.config.journal.Record(stepSuccessful{Id: stepContext.Id})
	}

	build.config.journal.Record(taskSuccessful{Id: task.Name})
	return nil
}

func buildStepContexts(logger *log.Logger, config *BuildConfig, task *project.Task) ([]StepContext, error) {
	var contexts []StepContext

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
		stepContext.Env.Secrets = config.Secrets

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
