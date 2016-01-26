package git

import (
	"io"
)

type cloneOptions struct {
	Path          string
	RepositoryUrl string
	Checkout      bool
	Branch        string
}

func NewCloneOptions(repositoryUrl string, path string) *cloneOptions {
	return &cloneOptions{
		Path:          path,
		RepositoryUrl: repositoryUrl,
		Branch:        "master",
		Checkout:      true,
	}
}

func (o cloneOptions) GetArgs() []string {
	args := []string{}

	if !o.Checkout {
		args = append(args, "--no-checkout")
	}

	args = append(args, []string{"--branch", o.Branch}...)
	args = append(args, o.RepositoryUrl)
	args = append(args, o.Path)
	return args
}

func (o cloneOptions) GetName() string {
	return "clone"
}

func (o cloneOptions) GetRepoPath() string {
	return o.Path
}

func Clone(options *cloneOptions, writer io.Writer) error {
	err := executeGit(options, writer)

	return err
}
