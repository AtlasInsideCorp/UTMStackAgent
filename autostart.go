package main

import (
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
		_, err := execute("nssm.exe", path, "install", "utmstack", "utmstack-windows.exe", "run")
		if err != nil {
			return err
		}

		execute("nssm.exe", path, "set", "utmstack", "AppDirectory", path)
		if err != nil {
			return err
		}

		execute("nssm.exe", path, "set", "utmstack", "DisplayName", "UTMStack Agent")
		if err != nil {
			return err
		}

		execute("nssm.exe", path, "set", "utmstack", "AppExit", "Default", "Restart")
		if err != nil {
			return err
		}

		execute("nssm.exe", path, "set", "utmstack", "Start", "SERVICE_AUTO_START")
		if err != nil {
			return err
		}

		execute("nssm.exe", path, "set", "utmstack", "ObjectName", "LocalSystem")
		if err != nil {
			return err
		}

		execute("nssm.exe", path, "start", "utmstack")
		if err != nil {
			return err
		}
	case "linux":
		type bash struct {
			Path string
		}

		scriptFile := filepath.Join("/", "usr", "local", "bin", "utmstack-agent.sh")
		scriptTemplateFile := filepath.Join(path, "templates", "utmstack-agent-bash.template")

		err := generateFromTemplate(bash{Path: path}, scriptTemplateFile, scriptFile)
		if err != nil {
			return err
		}

		_, err = execute("chmod", path, "755", filepath.Join("/", "usr", "local", "bin", "utmstack-agent.sh"))
		if err != nil {
			return err
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

		_, err = execute("systemctl", filepath.Join(path, "beats"), "enable", "utmstack")
		if err != nil {
			return err
		}

		_, err = execute("systemctl", filepath.Join(path, "beats"), "restart", "utmstack")
		if err != nil {
			return err
		}
	}
	return nil
}
