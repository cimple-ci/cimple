package project

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/mitchellh/mapstructure"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type publishDestination interface {
	Execute(files []string, vars StepVars, stdout io.Writer, stderr io.Writer) error
}

type PublishParser struct {
}

func (p PublishParser) GetToken() string {
	return "publish"
}

func (p PublishParser) Parse(item *ast.ObjectItem) (Step, error) {
	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return nil, err
	}

	name := item.Keys[0].Token.Value().(string)

	var c PublishStep
	c.name = name
	c.env = make(map[string]string)
	c.Destinations = make([]publishDestination, 0)

	a := item.Val.(*ast.ObjectType).List
	destinations := a.Filter("destination")
	for _, d := range destinations.Items {
		destination, err := parseDestination(d)
		if err != nil {
			return nil, err
		}
		c.Destinations = append(c.Destinations, destination)
	}

	delete(m, "env")

	if err := mapstructure.WeakDecode(m, &c); err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := parseEnvs(c.env, a.Filter("env")); err != nil {
		return nil, err
	}

	return c, nil
}

func parseDestination(item *ast.ObjectItem) (publishDestination, error) {
	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		log.Fatal(err)
		return nil, err
	}
	destination := &bintrayPublishDestination{}

	if err := mapstructure.WeakDecode(m, &destination); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return destination, nil
}

type PublishStep struct {
	name         string
	Files        []string
	Skip         bool
	Destinations []publishDestination
	env          map[string]string
}

func (c PublishStep) GetName() string {
	return c.name
}

func (c PublishStep) GetSkip() bool {
	return c.Skip
}

func (c PublishStep) GetEnv() map[string]string {
	return c.env
}

func (c PublishStep) Execute(vars StepVars, stdout io.Writer, stderr io.Writer) error {
	for _, destination := range c.Destinations {
		files := []string{}
		for _, f := range c.Files {
			path, err := templateString(f, vars)
			if err != nil {
				return err
			}
			matches, err := filepath.Glob(path)
			if err != nil {
				return err
			}
			files = append(files, matches...)
		}
		err := destination.Execute(files, vars, stdout, stderr)
		if err != nil {
			return err
		}
	}

	return nil
}

func templateString(s string, vars StepVars) (string, error) {
	tmpl, err := template.New("t").Parse(s)
	if err != nil {
		return "", err
	}
	var doc bytes.Buffer
	tmpl.Execute(&doc, vars)
	res := doc.String()

	return res, nil
}

type bintrayPublishDestination struct {
	Subject    string
	Repository string
	Package    string
	Username   string
}

func (b bintrayPublishDestination) Execute(files []string, vars StepVars, stdout io.Writer, stderr io.Writer) error {
	subject := b.Subject
	repository := b.Repository
	version := vars.Project.Version
	pkg := b.Package
	username := b.Username
	password, err := vars.Secrets.Get("bintray", b.Username)
	if err != nil {
		return err
	}

	for _, f := range files {
		filePath, err := filepath.Abs(f)
		if err != nil {
			return err
		}
		fileName := filepath.Base(filePath)
		url := fmt.Sprintf("https://api.bintray.com/content/%s/%s/%s/%s/%s", subject, repository, pkg, version, fileName)
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		fi, err := file.Stat()
		if err != nil {
			return err
		}

		req, err := http.NewRequest("PUT", url, file)
		if err != nil {
			return err
		}
		req.ContentLength = int64(fi.Size())
		req.SetBasicAuth(username, password)

		client := &http.Client{}

		log.Printf("Publishing %s to %s", filePath, url)

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)

		if res.StatusCode != 201 {
			return fmt.Errorf("Failed to publish %s. Receieved %d response - %s", filePath, res.StatusCode, body)
		} else {
			log.Printf("Published %s", filePath)
		}
	}

	return nil
}
