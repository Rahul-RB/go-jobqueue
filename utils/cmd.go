package utils

import (
	"errors"
	"io"
	"os/exec"
)

type Cmd struct {
	Cmd    *exec.Cmd
	Stdout *io.ReadCloser
	Stderr *io.ReadCloser
}

func RunWithTimeout(path string, args ...string) (*Cmd, error) {
	cmd := exec.Command(path, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return &Cmd{
			Cmd:    cmd,
			Stdout: &stdout,
			Stderr: nil,
		}, errors.New("Failed to pipe stdout: " + err.Error())

	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return &Cmd{
			Cmd:    cmd,
			Stdout: &stdout,
			Stderr: &stderr,
		}, errors.New("Failed to pipe stderr: " + err.Error())
	}

	err = cmd.Start()
	if err != nil {
		return &Cmd{
			Cmd:    cmd,
			Stdout: &stdout,
			Stderr: &stderr,
		}, errors.New("Failed to start command: " + err.Error())
	}

	return &Cmd{
		Cmd:    cmd,
		Stdout: &stdout,
		Stderr: &stderr,
	}, nil
}
