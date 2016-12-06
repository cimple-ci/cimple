package project

import (
	"log"

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
