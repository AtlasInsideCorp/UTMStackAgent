package main

import (
	"fmt"
	"os/exec"
)

func detectLinuxFamily() (string, error) {
	var pmCommands map[string]string = map[string]string{
		"debian": "dpkg-architecture",
		"rhel":   "yum",
	}

	for dist, command := range pmCommands {
		_, err := exec.Command(command).Output()
		if err == nil {
			return dist, nil
		}
	}
	return "", fmt.Errorf("unknown distribution")
}
