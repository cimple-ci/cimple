package vcs

import (
	"bytes"
	"fmt"
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
	branch, err := executeGit("rev-parse --abbrev-ref HEAD")
	if err != nil {
		return err
	}

	info.Branch = strings.TrimRight(string(branch), "\n")
	return nil
}

func remoteUrl(info *VcsInformation) error {
	url, err := executeGit(fmt.Sprintf("config --get remote.%s.url", info.RemoteName))
	if err != nil {
		return nil
	}

	info.RemoteUrl = string(url)
	return nil
}

func remoteName(info *VcsInformation) error {
	remote, err := executeGit(fmt.Sprintf("config --get branch.%s.remote", info.Branch))
	if err != nil {
		return nil
	}

	info.RemoteName = string(remote)
	return nil
}

func currentHash(info *VcsInformation) error {
	hash, err := executeGit("log -n 1 --pretty=format:%H")
	if err != nil {
		return err
	}

	info.Revision = string(hash)
	return nil
}

func executeGit(arguments string) ([]byte, error) {
	args := strings.Fields(arguments)
	var cmd = exec.Command("git", args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	out := buf.Bytes()
	if err == nil {
		return out, nil
	}

	return out, err
}
