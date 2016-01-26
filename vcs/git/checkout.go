package git

import (
	"io"
)

type checkoutOptions struct {
	Path   string
	Branch string
}

func NewCheckoutOptions(path string, branch string) *checkoutOptions {
	return &checkoutOptions{
		Path:   path,
		Branch: branch,
	}
}

func (o checkoutOptions) GetArgs() []string {
	args := []string{}

	args = append(args, o.Branch)
	return args
}

func (o checkoutOptions) GetName() string {
	return "checkout"
}

func (o checkoutOptions) GetRepoPath() string {
	return o.Path
}

func Checkout(options *checkoutOptions, writer io.Writer) error {
	err := executeGit(options, writer)

	return err
}
