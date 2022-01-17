package main

import (
	"os/exec"
)

func execute(c string, dir string, arg ...string) (string, error) {
	cmd := exec.Command(c, arg...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		h.Debug("Command message: %s", string(out[:]))
		h.Debug("Command error: %s", err)
		return "", err
	}

	return string(out[:]), nil
}
