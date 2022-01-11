package main

import (
	"os/exec"
)

func execute(c string, dir string, arg ...string) (string, error) {
	var err error
	if out, err := exec.Command(c, arg...).Output(); err == nil {
		output := string(out[:])
		return output, nil
	}
	return "", err
}
