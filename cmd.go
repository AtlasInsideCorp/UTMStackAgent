package main

import (
	"os/exec"
)

func execute(c string, dir string, arg ...string) (string, bool) {
	cmd := exec.Command(c, arg...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return string(out[:]) + err.Error(), true
	}

	return string(out[:]), false
}
