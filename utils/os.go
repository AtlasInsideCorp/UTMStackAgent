package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// DetectLinuxFamily is a function that detects the Linux distribution the program is running on.
// Returns the name of the detected distribution and an error if the distribution cannot be determined.
func DetectLinuxFamily() (string, error) {
	var pmCommands map[string]string = map[string]string{
		"debian": "apt list",
		"rhel":   "yum list",
	}

	for dist, command := range pmCommands {
		cmd := strings.Split(command, " ")
		var err error

		if len(cmd) > 1 {
			_, err = exec.Command(cmd[0], cmd[1:]...).Output()
		} else {
			_, err = exec.Command(cmd[0]).Output()
		}

		if err == nil {
			return dist, nil
		}
	}
	return "", fmt.Errorf("unknown distribution")
}
