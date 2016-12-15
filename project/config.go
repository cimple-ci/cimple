package project

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/lukesmith/cimple/env"
	"github.com/lukesmith/cimple/vcs"
	"github.com/mitchellh/mapstructure"
	"regexp"
	"strings"
	"time"
)

type SecretStore interface {
	Get(t string, k string) (string, error)
}

type StepParser interface {
	GetToken() string
	Parse(*ast.ObjectItem) (Step, error)
}

type Task struct {
	Description string
	Depends     []string
	Name        string
	Steps       map[string]Step
	StepOrder   []string
	Archive     []string
	Env         map[string]string
	Skip        bool
	LimitTo     string
}

func (t Task) GetID() string {
	return t.Name
}

func (t Task) GetDependencies() []string {
	return t.Depends
}

type StepVars struct {
	Cimple     *env.CimpleEnvironment
	BuildDate  time.Time
	Project    Project
	Vcs        vcs.VcsInformation
	TaskName   string
	WorkingDir string
	HostEnv    map[string]string
	StepEnv    map[string]string
	Secrets    SecretStore
}

func (sv StepVars) FormattedBuildDate() string {
	return sv.BuildDate.Format(time.RFC3339)
}

func (sv *StepVars) Map() map[string]string {
	m := make(map[string]string)
	m = merge(m, sv.HostEnv)

	m["CIMPLE_BUILD_DATE"] = sv.BuildDate.Format(time.RFC3339)
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

type Step interface {
	GetSkip() bool
	GetName() string
	GetEnv() map[string]string
	Execute(vars StepVars, stdout io.Writer, stderr io.Writer) error
}

type Config struct {
	Project Project
	Tasks   map[string]*Task
}

type Project struct {
	Name        string
	Description string
	Version     string
	Env         map[string]string
}

type ConfigError struct {
	Issues []string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("One or more configuration errors exist:\n%s", strings.Join(e.Issues, "\n"))
}

func LoadConfig(path string) (*Config, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Load(string(d))
}

func Load(str string) (*Config, error) {
	obj, err := hcl.Parse(str)
	if err != nil {
		return nil, err
	}

	cfg, err := parseConfig(obj)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseConfig(obj *ast.File) (*Config, error) {
	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, obj); err != nil {
		return nil, err
	}

	var result Config
	result.Tasks = make(map[string]*Task)

	if err := mapstructure.WeakDecode(m, &result); err != nil {
		return nil, err
	}

	for _, fieldName := range []string{"name", "version"} {
		if _, ok := m[fieldName]; !ok {
			return nil, fmt.Errorf("'%s' was not specified", fieldName)
		}
	}

	result.Project.Name = m["name"].(string)
	result.Project.Version = m["version"].(string)

	if val, ok := m["description"]; ok {
		result.Project.Description = val.(string)
	}

	result.Project.Env = make(map[string]string)

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("Node is not an ObjectList")
	}

	matches := list.Filter("task")
	for _, m := range matches.Items {
		err := parseTask(result.Tasks, m)
		if err != nil {
			return nil, err
		}
	}

	if err := parseEnvs(result.Project.Env, list.Filter("env")); err != nil {
		return nil, err
	}

	return &result, nil
}

func parseTask(tasks map[string]*Task, item *ast.ObjectItem) error {
	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return err
	}

	delete(m, "env")

	var task Task
	task.Name = item.Keys[0].Token.Value().(string)
	task.Env = make(map[string]string)
	task.Steps = make(map[string]Step)
	task.StepOrder = []string{}
	task.Depends = []string{}
	task.Archive = []string{}

	if err := mapstructure.WeakDecode(m, &task); err != nil {
		return err
	}

	if val, ok := m["limit_to"]; ok {
		task.LimitTo = val.(string)
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	}

	stepParsers := []StepParser{&ScriptStepParser{}, &CommandStepParser{}, &PublishParser{}}

	so, err := stepOrder(listVal)
	if err != nil {
		return err
	}
	task.StepOrder = so

	for _, sp := range stepParsers {
		if o := listVal.Filter(sp.GetToken()); len(o.Items) > 0 {
			steps := make([]Step, 0)

			for _, item := range o.Items {
				step, err := sp.Parse(item)
				if err != nil {
					return err
				}
				steps = append(steps, step)
			}

			for _, st := range steps {
				if _, exists := task.Steps[st.GetName()]; exists {
					return &ConfigError{
						Issues: []string{fmt.Sprintf("A step named %s exists multiple times", st.GetName())},
					}
				}
				task.Steps[st.GetName()] = st
			}
		}
	}

	if err := parseEnvs(task.Env, listVal.Filter("env")); err != nil {
		return err
	}

	_, exists := tasks[task.Name]
	if exists {
		return &ConfigError{
			Issues: []string{fmt.Sprintf("A task named %s exists multiple times", task.Name)},
		}
	}

	tasks[task.Name] = &task

	return nil
}

func stepOrder(o *ast.ObjectList) ([]string, error) {
	result := []string{}

	for _, item := range o.Items {
		for _, keyItem := range item.Keys {
			key := keyItem.Token.Value().(string)

			if key == "command" || key == "script" || key == "publish" {
				n := item.Keys[1].Token.Value().(string)
				result = append(result, n)
			}
		}
	}

	return result, nil
}

func validateConfig(cfg *Config) error {
	var issues = []string{}

	r, _ := regexp.Compile("[a-z0-9_]+")

	for taskName, task := range cfg.Tasks {
		if !r.MatchString(taskName) {
			issues = append(issues, fmt.Sprintf("%s is not a valid task name", taskName))
		}

		for stepName, _ := range task.Steps {
			if !r.MatchString(stepName) {
				issues = append(issues, fmt.Sprintf("%s is not a valid step name", stepName))
			}
		}
	}

	if len(issues) > 0 {
		return &ConfigError{
			Issues: issues,
		}
	}

	return nil
}

func parseEnvs(result map[string]string, list *ast.ObjectList) error {
	for _, item := range list.Elem().Items {
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			log.Fatal(err)
			return err
		}

		if err := mapstructure.WeakDecode(m, &result); err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

func count(s []string, e string) int {
	var occurrences = 0
	for _, a := range s {
		if a == e {
			occurrences = occurrences + 1
		}
	}
	return occurrences
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
