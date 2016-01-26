package git

import (
	"testing"
)

func TestCloneGetName(t *testing.T) {
	options := &cloneOptions{}

	if options.GetName() != "clone" {
		t.Errorf("Expected name of command to be clone, was %s", options.GetName())
	}
}

func TestCloneGetRepoPath(t *testing.T) {
	options := &cloneOptions{
		Path: "/tmp/repo_path",
	}

	if options.GetRepoPath() != options.Path {
		t.Errorf("Expected repo path to return %s, was %s", options.Path, options.GetRepoPath())
	}
}

func TestCloneGetArgsWithDefaults(t *testing.T) {
	options := NewCloneOptions("git://github.com/test.git", "/tmp/test_path")

	args := options.GetArgs()

	if args[0] != "--branch" {
		t.Fatalf("Expected argument to be the branch flag, was %s", args[0])
	}

	if args[1] != options.Branch {
		t.Fatalf("Expected argument to be the branch name, was %s", args[1])
	}

	if args[2] != options.RepositoryUrl {
		t.Fatalf("Expected argument to be the repository url, was %s", args[2])
	}

	if args[3] != options.Path {
		t.Fatalf("Expected argument to be the path, was %s", args[3])
	}
}

func TestCloneGetArgsWithNoCheckout(t *testing.T) {
	options := NewCloneOptions("git://github.com/test.git", "/tmp/test_path")
	options.Checkout = false

	args := options.GetArgs()

	if args[0] != "--no-checkout" {
		t.Fatal("Expected --no-checkout to be present")
	}
}
