package vcs

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type VcsInformation struct {
	Vcs        string
	Branch     string
	Revision   string
	RemoteUrl  string
	RemoteName string
}

func LoadVcsInformation() (*VcsInformation, error) {
	info := new(VcsInformation)
	info.Vcs = "Git"

	err := currentBranch(info)
	if err != nil {
		return nil, err
	}

	err = currentHash(info)
	if err != nil {
		return nil, err
	}

	err = remoteName(info)
	if err != nil {
		return nil, err
	}

	err = remoteUrl(info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func currentBranch(info *VcsInformation) error {
	buf := &bytes.Buffer{}
	err := executeGit("rev-parse --abbrev-ref HEAD", buf)
	if err != nil {
		return err
	}

	branch := buf.Bytes()

	info.Branch = strings.TrimRight(string(branch), "\n")
	return nil
}

func remoteUrl(info *VcsInformation) error {
	buf := &bytes.Buffer{}
	err := executeGit(fmt.Sprintf("config --get remote.%s.url", info.RemoteName), buf)
	if err != nil {
		return nil
	}

	info.RemoteUrl = string(buf.Bytes())
	return nil
}

func remoteName(info *VcsInformation) error {
	buf := &bytes.Buffer{}
	err := executeGit(fmt.Sprintf("config --get branch.%s.remote", info.Branch), buf)
	if err != nil {
		return nil
	}

	info.RemoteName = string(buf.Bytes())
	return nil
}

func currentHash(info *VcsInformation) error {
	buf := &bytes.Buffer{}
	err := executeGit("log -n 1 --pretty=format:%H", buf)
	if err != nil {
		return err
	}

	info.Revision = string(buf.Bytes())
	return nil
}

func executeGit(arguments string, writer io.Writer) error {
	args := strings.Fields(arguments)
	var cmd = exec.Command("git", args...)
	cmd.Stdout = writer
	cmd.Stderr = writer

	err := cmd.Run()
	if err == nil {
		return nil
	}

	return err
}
