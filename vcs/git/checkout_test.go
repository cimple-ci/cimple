package git

import (
	"testing"
)

func TestCheckoutGetName(t *testing.T) {
	options := &checkoutOptions{}

	if options.GetName() != "checkout" {
		t.Errorf("Expected name of command to be checkout, was %s", options.GetName())
	}
}

func TestCheckoutGetRepoPath(t *testing.T) {
	options := &checkoutOptions{
		Path: "/tmp/repo_path",
	}

	if options.GetRepoPath() != options.Path {
		t.Errorf("Expected repo path to return %s, was %s", options.Path, options.GetRepoPath())
	}
}

func TestCheckoutGetArgsWithDefaults(t *testing.T) {
	options := NewCheckoutOptions("/tmp/test_path", "the-branch")

	args := options.GetArgs()

	if args[0] != options.Branch {
		t.Fatalf("Expected argument to be the branch name, was %s", args[0])
	}
}
