package utils

import (
	"errors"
	"os/exec"
)

// Execute runs the specified command in the specified directory .
// Returns console output and a boolean indicating whether an error occurred
func Execute(c string, dir string, arg ...string) (string, bool) {
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
