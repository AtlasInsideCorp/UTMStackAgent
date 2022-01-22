package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func detectLinuxFamily() (string, error) {
	var pmCommands map[string]string = map[string]string{
		"debian": "apt list",
		"rhel":   "yum",
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
