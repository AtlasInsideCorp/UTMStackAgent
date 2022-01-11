package main

import "runtime"

func autoStart() error {
	path, err := getMyPath()
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "windows":
		_, err := execute("nssm.exe", path, "install", "utmagent", "utmagent.exe")
		if err != nil {
			return err
		}
		execute("nssm.exe", path, "set", "utmagent", "AppDirectory", path)
		if err != nil {
			return err
		}
		execute("nssm.exe", path, "set", "utmagent", "DisplayName", "UTMStack Agent")
		if err != nil {
			return err
		}
		execute("nssm.exe", path, "set", "utmagent", "AppExit", "Default", "Restart")
		if err != nil {
			return err
		}
		execute("nssm.exe", path, "set", "utmagent", "Start", "SERVICE_AUTO_START")
		if err != nil {
			return err
		}
		execute("nssm.exe", path, "set", "utmagent", "ObjectName", "LocalSystem")
		if err != nil {
			return err
		}
		execute("nssm.exe", path, "start", "utmagent")
		if err != nil {
			return err
		}
	}
	return nil
}
