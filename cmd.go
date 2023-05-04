package main

import (
	"errors"
	"os/exec"
)

func execute(c string, dir string, arg ...string) (string, bool) {
	cmd := exec.Command(c, arg...)

	cmd.Dir = dir
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}

	out, err := cmd.Output()
	if err != nil {
		return string(out[:]) + err.Error(), true
	}

	return string(out[:]), false
}
