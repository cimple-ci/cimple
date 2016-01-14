package database

import (
	"fmt"
	"path/filepath"
)

type CimpleDatabase interface {
	GetProjects() []*Project
	GetProject(name string) (*Project, error)
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

func NewDatabase(path string) CimpleDatabase {
	return &database{
		path: path,
	}
}

type Project struct {
	Name       string
	BuildCount int
}
