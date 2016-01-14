package database

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

type CimpleDatabase interface {
	GetProjects() []*Project
	GetProject(name string) (*Project, error)
	GetBuilds(project string) ([]*Build, error)
	GetBuild(project string, id string) (*Build, error)
}

type database struct {
	path string
}

func (db *database) GetProjects() []*Project {
	dirs, _ := filepath.Glob(filepath.Join(db.path, "*"))

	projects := []*Project{}

	for _, d := range dirs {
		builds, _ := filepath.Glob(filepath.Join(d, "*"))
		projects = append(projects, &Project{
			Name:       filepath.Base(d),
			BuildCount: len(builds),
		})
	}

	return projects
}

func (db *database) GetProject(name string) (*Project, error) {
	for _, p := range db.GetProjects() {
		if p.Name == name {
			return p, nil
		}
	}

	return nil, fmt.Errorf("Unable to find project named %s", name)
}

func (db *database) GetBuilds(project string) ([]*Build, error) {
	dirs, _ := filepath.Glob(filepath.Join(db.path, project, "*"))

	builds := []*Build{}

	for _, d := range dirs {
		t, err := time.Parse(time.RFC3339, filepath.Base(d))
		if err != nil {
			return []*Build{}, err
		}

		builds = append(builds, &Build{
			Id:         filepath.Base(d),
			Date:       t,
			outputPath: filepath.Join(d, "output"),
		})
	}

	return builds, nil
}

func (db *database) GetBuild(project string, id string) (*Build, error) {
	builds, err := db.GetBuilds(project)
	if err != nil {
		return nil, err
	}

	for _, p := range builds {
		if p.Id == id {
			return p, nil
		}
	}

	return nil, fmt.Errorf("Unable to find build %s for project %s", id, project)
}

func NewDatabase(path string) CimpleDatabase {
	return &database{
		path: path,
	}
}

type Project struct {
	Name       string
	BuildCount int
}

type Build struct {
	Id         string
	Date       time.Time
	outputPath string
}

func (b *Build) GetOutput() ([]byte, error) {
	return ioutil.ReadFile(b.outputPath)
}
