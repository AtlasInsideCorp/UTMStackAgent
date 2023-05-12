package utils

import (
	"fmt"
	"runtime"
)

func StopService(name string) error {
	path, err := GetMyPath()
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "windows":
		result, err := Execute("cmd", path, "/c", "sc", "stop", name)
		if err {
			return fmt.Errorf("error stoping service: %v", result)
		}
	case "linux":
		result, err := Execute("systemctl", path, "stop", name)
		if err {
			return fmt.Errorf("error stoping service: %v", result)
		}
	}
	return nil
}

func UninstallService(name string) error {
	path, err := GetMyPath()
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "windows":
		result, err := Execute("cmd", path, "/c", "sc", "delete", name)
		if err {
			return fmt.Errorf("error uninstalling service: %v", result)
		}
	case "linux":
		result, err := Execute("systemctl", path, "disable", name)
		if err {
			return fmt.Errorf("error uninstalling service: %v", result)
		}
		result, err = Execute("rm", "/etc/systemd/system/"+name+".service", "/etc/systemd/system/"+name+".service")
		if err {
			return fmt.Errorf("error uninstalling service: %v", result)
		}
	}
	return nil
}
