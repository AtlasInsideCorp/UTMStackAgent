package main

import "runtime"

func uninstall() error {
	path, err := getMyPath()
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "windows":
		_, err := execute("nssm.exe", path, "stop", "utmagent")
		if err != nil {
			return err
		}
		_, err = execute("nssm.exe", path, "remove", "utmagent")
		if err != nil {
			return err
		}
	}
	return nil
}
