package main

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func uninstall() error {
	path, err := getMyPath()
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "windows":
		result, err := execute("nssm.exe", path, "stop", "utmstack")
		if err{
			return fmt.Errorf("%s", result)
		}

		result, err = execute("nssm.exe", path, "remove", "utmstack", "confirm")
		if err{
			return fmt.Errorf("%s", result)
		}
	case "linux":
		result, err := execute("systemctl", path, "disable", "utmstack")
		if err{
			return fmt.Errorf("%s", result)
		}

		result, err = execute("systemctl", path, "stop", "utmstack")
		if err{
			return fmt.Errorf("%s", result)
		}

		result, err = execute("rm", path, filepath.Join("/", "etc", "systemd", "system", "utmstack.service"))
		if err{
			return fmt.Errorf("%s", result)
		}

		result, err = execute("systemctl", path, "daemon-reload")
		if err{
			return fmt.Errorf("%s", result)
		}
	}
	return nil
}
