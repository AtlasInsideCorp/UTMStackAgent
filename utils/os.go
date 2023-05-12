package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

// GetOSInfo gets information about the operating system the application is running on.
func GetOSInfo() (string, string, string) {
	var osName, osVersion, osPlatform string

	switch runtime.GOOS {
	case "linux":
		osName = "Linux"
		file, err := os.ReadFile("/etc/os-release")
		if err != nil {
			osVersion = "unknown"
		} else {
			lines := strings.Split(string(file), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "VERSION=") {
					osVersion = strings.Trim(strings.TrimPrefix(line, "VERSION="), `"`)
					break
				}
			}
		}
		osPlatform = runtime.GOARCH
	//case "darwin":
	//	osName = "macOS"
	//	osVersion = "unknown"
	//	osPlatform = runtime.GOARCH
	case "windows":
		osName = "Windows"
		v, err := exec.Command("ver").Output()
		if err != nil {
			osVersion = "unknown"
		} else {
			osVersion = string(v)
		}
		osPlatform = runtime.GOARCH
	default:
		osName = "unknown"
		osVersion = "unknown"
		osPlatform = "unknown"
	}

	return osName, osVersion, osPlatform
}
