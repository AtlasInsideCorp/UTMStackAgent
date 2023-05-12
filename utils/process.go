package utils

import (
	"os/exec"
	"runtime"
	"strings"
)

func IsProcessRunning(processName string) (bool, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist")
	} else {
		cmd = exec.Command("ps", "aux")
	}

	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	processes := strings.Split(string(output), "\n")
	for _, process := range processes {
		if strings.Contains(process, processName) {
			return true, nil
		}
	}

	return false, nil
}
