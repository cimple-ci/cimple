package git

import (
	"io"
	"log"
	"os/exec"
	"github.com/lukesmith/cimple/logging"
)

type GitCommand interface {
	GetName() string
	GetArgs() []string
	GetRepoPath() string
}

func executeGit(command GitCommand, writer io.Writer) error {
	logger := logging.CreateLogger("Git", writer)
	logger.Printf("Performing git %s %s", command.GetName(), command.GetArgs())

	args := []string{command.GetName()}
	args = append(args, command.GetArgs()...)
	var cmd = exec.Command("git", args...)
	cmd.Dir = command.GetRepoPath()
	cmd.Stdout = writer
	cmd.Stderr = writer

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
