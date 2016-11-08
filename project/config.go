package project

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/mitchellh/mapstructure"
	"regexp"
	"strings"
)

type Task struct {
	Description string
	Name        string
	Steps       map[string]Step
	StepOrder   []string
	Archive     []string
	Env         map[string]string
	Skip        bool
}

type Step interface {
	GetSkip() bool
	GetEnv() map[string]string
}

type Script struct {
	Skip bool
	Body string
	Env  map[string]string
}

func (s Script) GetSkip() bool {
	return s.Skip
}

func (s Script) GetEnv() map[string]string {
	return s.Env
}

type Command struct {
	Command string
	Args    []string
	Env     map[string]string
	Skip    bool
}

func (c Command) GetSkip() bool {
	return c.Skip
}

func (c Command) GetEnv() map[string]string {
	return c.Env
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

	if o := list.Filter("env"); len(o.Items) > 0 {
		if err := parseEnvs(result.Project.Env, o); err != nil {
			return nil, err
		}
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

	if err := mapstructure.WeakDecode(m, &task); err != nil {
		return err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	}

	so, err := stepOrder(listVal)
	if err != nil {
		return err
	}
	task.StepOrder = so

	if o := listVal.Filter("command"); len(o.Items) > 0 {
		if err := parseCommands(task.Steps, o); err != nil {
			return err
		}
	}

	if o := listVal.Filter("script"); len(o.Items) > 0 {
		if err := parseScripts(task.Steps, o); err != nil {
			return err
		}
	}

	if o := listVal.Filter("env"); len(o.Items) > 0 {
		if err := parseEnvs(task.Env, o); err != nil {
			return err
		}
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
	var result []string

	for _, item := range o.Items {
		for _, keyItem := range item.Keys {
			key := keyItem.Token.Value().(string)

			if key == "command" || key == "script" {
				n := item.Keys[1].Token.Value().(string)
				result = append(result, n)

				if count(result, n) > 1 {
					return nil, &ConfigError{
						Issues: []string{fmt.Sprintf("A step named %s exists multiple times", n)},
					}
				}
			}
		}
	}

	return result, nil
}

func parseCommands(steps map[string]Step, list *ast.ObjectList) error {
	for _, item := range list.Items {
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return err
		}

		delete(m, "env")

		name := item.Keys[0].Token.Value().(string)
		var c Command
		c.Env = make(map[string]string)
		if err := mapstructure.WeakDecode(m, &c); err != nil {
			log.Fatal(err)
			return err
		}

		steps[name] = c

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		}

		if o := listVal.Filter("env"); len(o.Items) > 0 {
			if err := parseEnvs(c.Env, o); err != nil {
				return err
			}
		}
	}

	return nil
}

func parseScripts(result map[string]Step, list *ast.ObjectList) error {
	for _, item := range list.Items {
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return err
		}

		delete(m, "env")

		name := item.Keys[0].Token.Value().(string)
		var c Script
		c.Env = make(map[string]string)
		if err := mapstructure.WeakDecode(m, &c); err != nil {
			log.Fatal(err)
			return err
		}
		result[name] = c

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		}

		if o := listVal.Filter("env"); len(o.Items) > 0 {
			if err := parseEnvs(c.Env, o); err != nil {
				return err
			}
		}
	}

	return nil
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
