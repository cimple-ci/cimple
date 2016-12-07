package project

import (
	"bytes"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/mitchellh/mapstructure"
	"io"
	"log"
	exec "os/exec"
	"text/template"
)

type CommandStepParser struct {
}

func (st CommandStepParser) GetToken() string {
	return "command"
}

func (st CommandStepParser) Parse(list *ast.ObjectList) ([]Step, error) {
	result := make([]Step, 0)

	for _, item := range list.Items {
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return nil, err
		}

		delete(m, "env")

		name := item.Keys[0].Token.Value().(string)
		var c Command
		c.Env = make(map[string]string)
		if err := mapstructure.WeakDecode(m, &c); err != nil {
			log.Fatal(err)
			return nil, err
		}

		c.name = name
		result = append(result, c)

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		}

		if o := listVal.Filter("env"); len(o.Items) > 0 {
			if err := parseEnvs(c.Env, o); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

type Command struct {
	name    string
	Command string
	Args    []string
	Env     map[string]string
	Skip    bool
}

func (c Command) GetName() string {
	return c.name
}

func (c Command) GetSkip() bool {
	return c.Skip
}

func (c Command) GetEnv() map[string]string {
	return c.Env
}

func (c Command) Execute(vars StepVars, stdout io.Writer, stderr io.Writer) error {
	args, err := c.templateArgs(vars)
	if err != nil {
		return err
	}

	var cmd = exec.Command(c.Command, args...)

	// Clear out env for command so not to inherit current process's environment.
	cmd.Env = []string{}

	env, err := c.templatedEnvs(vars)
	if err != nil {
		return err
	}

	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (c Command) templateArgs(vars StepVars) ([]string, error) {
	args := []string{}
	for _, v := range c.Args {
		tmpl, err := template.New("t").Parse(v)
		if err != nil {
			return nil, err
		}
		var doc bytes.Buffer
		tmpl.Execute(&doc, vars)
		args = append(args, doc.String())
	}

	return args, nil
}

func (c Command) templatedEnvs(vars StepVars) (map[string]string, error) {
	env := make(map[string]string)

	for k, v := range vars.Map() {
		tmpl, err := template.New("t").Parse(v)
		if err != nil {
			return nil, err
		}
		var doc bytes.Buffer
		tmpl.Execute(&doc, vars)
		env[k] = doc.String()
	}

	return env, nil
}
