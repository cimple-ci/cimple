package database

type CimpleDatabase interface {
	GetProjects() []*Project
}

type database struct {
}

func (db *database) GetProjects() []*Project {
	projects := []*Project{
		&Project{Name: "Cimple"},
		&Project{Name: "Cimple car"},
	}

	return projects
}

func NewDatabase() CimpleDatabase {
	return &database{}
}

type Project struct {
	Name string
}
