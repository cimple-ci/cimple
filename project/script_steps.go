package project

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	exec "os/exec"
	"text/template"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/mitchellh/mapstructure"
)

type ScriptStepParser struct {
}

func (st ScriptStepParser) GetToken() string {
	return "script"
}

func (st ScriptStepParser) Parse(list *ast.ObjectList) ([]Step, error) {
	result := make([]Step, 0)

	for _, item := range list.Items {
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return nil, err
		}

		delete(m, "env")

		name := item.Keys[0].Token.Value().(string)
		var c Script
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

type Script struct {
	name string
	Skip bool
	Body string
	Env  map[string]string
}

func (s Script) GetName() string {
	return s.name
}

func (s Script) GetSkip() bool {
	return s.Skip
}

func (s Script) GetEnv() map[string]string {
	return s.Env
}

func (s Script) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Body string
		Skip bool
	}{
		s.Body,
		s.Skip,
	})
}

func (s Script) Execute(vars StepVars, stdout io.Writer, stderr io.Writer) error {
	f, err := s.writeFile(vars)
	if err != nil {
		return err
	}

	args := []string{f}
	var cmd = exec.Command("/bin/sh", args...)

	// Clear out env for command so not to inherit current process's environment.
	cmd.Env = []string{}

	env, err := s.templatedEnvs(vars)
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

func (s Script) templatedEnvs(vars StepVars) (map[string]string, error) {
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

func (s Script) writeFile(vars StepVars) (string, error) {
	f, _ := ioutil.TempFile(os.TempDir(), "step")
	defer f.Close()

	tmpl, err := template.New("t").Parse(s.Body)
	if err != nil {
		return "", err
	}
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, vars)
	if err != nil {
		return "", err
	}
	f.WriteString(doc.String())

	return f.Name(), nil
}
