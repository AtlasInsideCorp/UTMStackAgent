package main

import (
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
		_, err := execute("nssm.exe", path, "stop", "utmstack")
		if err != nil {
			return err
		}

		_, err = execute("nssm.exe", path, "remove", "utmstack")
		if err != nil {
			return err
		}
	case "linux":
		_, err = execute("systemctl", path, "disable", "utmstack")
		if err != nil {
			return err
		}

		_, err = execute("systemctl", path, "stop", "utmstack")
		if err != nil {
			return err
		}

		_, err = execute("rm", path, filepath.Join("/", "etc", "systemd", "system", "utmstack.service"))
		if err != nil {
			return err
		}

		_, err = execute("systemctl", path, "daemon-reload")
		if err != nil {
			return err
		}
	}
	return nil
}
