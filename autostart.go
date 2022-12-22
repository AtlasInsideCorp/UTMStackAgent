package main

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func autoStart() error {
	path, err := getMyPath()
	if err != nil {
		return err
	}
	switch runtime.GOOS {
	case "windows":
		result, err := execute(".\nssm.exe", path, "install", "utmstack", "utmstack-windows.exe", "run")
		if err {
			return fmt.Errorf("%s", result)
		}

		result, err = execute(".\nssm.exe", path, "set", "utmstack", "AppDirectory", path)
		if err {
			return fmt.Errorf("%s", result)
		}

		result, err = execute(".\nssm.exe", path, "set", "utmstack", "DisplayName", "UTMStack Agent")
		if err {
			return fmt.Errorf("%s", result)
		}

		result, err = execute(".\nssm.exe", path, "set", "utmstack", "AppExit", "Default", "Restart")
		if err {
			return fmt.Errorf("%s", result)
		}

		result, err = execute(".\nssm.exe", path, "set", "utmstack", "Start", "SERVICE_AUTO_START")
		if err {
			return fmt.Errorf("%s", result)
		}

		result, err = execute(".\nssm.exe", path, "set", "utmstack", "ObjectName", "LocalSystem")
		if err {
			return fmt.Errorf("%s", result)
		}

		result, err = execute(".\nssm.exe", path, "start", "utmstack")
		if err {
			return fmt.Errorf("%s", result)
		}
	case "linux":
		type bash struct {
			Path string
		}

		scriptFile := filepath.Join("/", "usr", "local", "bin", "utmstack-agent.sh")
		scriptTemplateFile := filepath.Join(path, "templates", "utmstack-agent-bash.sh")

		err := generateFromTemplate(bash{Path: path}, scriptTemplateFile, scriptFile)
		if err != nil {
			return err
		}

		result, errB := execute("chmod", path, "755", filepath.Join("/", "usr", "local", "bin", "utmstack-agent.sh"))
		if errB {
			return fmt.Errorf("%s", result)
		}

		incidentServiceConfig := `[Unit]
Description=UTMStack Agent
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=60
ExecStart=/usr/local/bin/utmstack-agent.sh
[Install]
WantedBy=multi-user.target`

		err = writeToFile("/etc/systemd/system/utmstack.service", incidentServiceConfig)
		if err != nil {
			return err
		}

		result, errB = execute("systemctl", filepath.Join(path, "beats"), "enable", "utmstack")
		if errB {
			return fmt.Errorf("%s", result)
		}

		result, errB = execute("systemctl", filepath.Join(path, "beats"), "restart", "utmstack")
		if errB {
			return fmt.Errorf("%s", result)
		}
	}
	return nil
}
