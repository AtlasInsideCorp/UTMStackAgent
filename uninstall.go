package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
)

func uninstall() error {
	path, err := utils.GetMyPath()
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "windows":
		result, err := utils.Execute("nssm.exe", path, "stop", "utmstack")
		if err {
			return fmt.Errorf("%s", result)
		}

		result, err = utils.Execute("nssm.exe", path, "remove", "utmstack", "confirm")
		if err {
			return fmt.Errorf("%s", result)
		}
	case "linux":
		result, err := utils.Execute("systemctl", path, "disable", "utmstack")
		if err {
			return fmt.Errorf("%s", result)
		}

		result, err = utils.Execute("systemctl", path, "stop", "utmstack")
		if err {
			return fmt.Errorf("%s", result)
		}

		result, err = utils.Execute("rm", path, filepath.Join("/", "etc", "systemd", "system", "utmstack.service"))
		if err {
			return fmt.Errorf("%s", result)
		}

		result, err = utils.Execute("systemctl", path, "daemon-reload")
		if err {
			return fmt.Errorf("%s", result)
		}
	}

	os.Remove(filepath.Join(path, "config.yml"))

	return nil
}
