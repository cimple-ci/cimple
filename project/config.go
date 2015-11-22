package project

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"log"
)

type Task struct {
	Description string
	Name        string
	Commands    map[string]Command
	Archive     []string
	Env         map[string]string
}

type Command struct {
	Command string
	Args    []string
	Env     map[string]string
	Skip    bool
}

type Config struct {
	Project Project
	Tasks   map[string]Task
}

type Project struct {
	Name        string
	Description string
}

func LoadConfig(path string) (*Config, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	obj, err := hcl.Parse(string(d))
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, obj); err != nil {
		return nil, err
	}

	var result Config
	result.Tasks = make(map[string]Task)

	if err := mapstructure.WeakDecode(m, &result); err != nil {
		return nil, err
	}

	result.Project.Name = m["name"].(string)
	result.Project.Description = m["description"].(string)

	list, ok := obj.Node.(*ast.ObjectList)
	if !ok {
		return nil, err
	}

	matches := list.Filter("task")
	for _, m := range matches.Items {
		err := parseTask(result.Tasks, m)
		if err != nil {
			return nil, nil
		}
	}

	return &result, nil
}

func parseTask(tasks map[string]Task, item *ast.ObjectItem) error {
	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return err
	}

	delete(m, "env")

	var task Task
	task.Name = item.Keys[0].Token.Value().(string)
	task.Env = make(map[string]string)
	task.Commands = make(map[string]Command)

	if err := mapstructure.WeakDecode(m, &task); err != nil {
		return err
	}

	var listVal *ast.ObjectList
	if ot, ok := item.Val.(*ast.ObjectType); ok {
		listVal = ot.List
	}

	if o := listVal.Filter("command"); len(o.Items) > 0 {
		if err := parseCommands(task.Commands, o); err != nil {
			return err
		}
	}

	if o := listVal.Filter("env"); len(o.Items) > 0 {
		if err := parseEnvs(task.Env, o); err != nil {
			return err
		}
	}

	tasks[task.Name] = task

	return nil
}

func parseCommands(result map[string]Command, list *ast.ObjectList) error {
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
