package main

import (
	"encoding/base64"
	"path/filepath"
	"runtime"
	"sync"
)

func startWazuh() {
	var runOnce sync.Once
	go func() {
		path, err := getMyPath()
		if err != nil {
			h.FatalError("error getting path: %v", err)
		}
		switch runtime.GOOS {
		case "windows":
			runOnce.Do(func() {
				_, err = execute(
					filepath.Join(path, "wazuh", "windows", "wazuh-agent.exe"),
					filepath.Join(path, "wazuh", "windows"),
					"start",
				)
				if err != nil {
					h.FatalError("error running wazuh: %v", err)
				}
			})
		}
	}()
}

func configureWazuh(ip, key string) error {
	path, err := getMyPath()
	if err != nil {
		return err
	}

	type WazuhConfig struct {
		IP string
	}

	dKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}

	config := WazuhConfig{ip}

	switch runtime.GOOS {
	case "windows":
		ossecFile := filepath.Join(path, "wazuh", "windows", "ossec.conf")
		ossecTemplateFile := filepath.Join(path, "templates", "wazuh-windows.conf")

		err := generateFromTemplate(config, ossecTemplateFile, ossecFile)
		if err != nil {
			return err
		}

		err = writeToFile(filepath.Join(path, "wazuh", "windows", "client.keys"), string(dKey[:]))
		if err != nil {
			return err
		}
	case "linux":
		var templateFile string
		configFile := filepath.Join("/", "var", "ossec", "etc", "ossec.conf")

		family, err := detectLinuxFamily()
		if err != nil {
			return err
		}

		switch family {
		case "debian":
			templateFile = filepath.Join(path, "templates", "wazuh-debian.conf")

			_, err := execute("apt", filepath.Join(path, "wazuh"), "update")
			if err != nil {
				return err
			}

			_, err = execute("apt", filepath.Join(path, "wazuh"), "install", "-y", "curl", "apt-transport-https", "lsb-release", "gnupg2", "wget")
			if err != nil {
				return err
			}

			err = download("https://packages.wazuh.com/key/GPG-KEY-WAZUH")
			if err != nil {
				return err
			}

			_, err = execute("apt-key", path, "add", "GPG-KEY-WAZUH")
			if err != nil {
				return err
			}

			err = writeToFile(filepath.Join("/", "etc", "apt", "sources.list.d", "wazuh.list"), "deb https://packages.wazuh.com/4.x/apt/ stable main")
			if err != nil {
				return err
			}

			_, err = execute("apt", filepath.Join(path, "wazuh"), "update")
			if err != nil {
				return err
			}

			_, err = execute("apt", filepath.Join(path, "wazuh"), "install", "-y", "wazuh-agent")
			if err != nil {
				return err
			}

		case "rhel":
			templateFile = filepath.Join(path, "templates", "wazuh-rhel.conf")

			_, err := execute("rpm", filepath.Join(path, "wazuh"), "--import", "https://packages.wazuh.com/key/GPG-KEY-WAZUH")
			if err != nil {
				return err
			}

			err = writeToFile(
				filepath.Join("/", "etc", "yum.repos.d", "wazuh.repo"),
				`[wazuh_repo]
gpgcheck=1
gpgkey=https://packages.wazuh.com/key/GPG-KEY-WAZUH
enabled=1
name=Wazuh repository
baseurl=https://packages.wazuh.com/4.x/yum/
protect=1`,
			)
			if err != nil {
				return err
			}

			_, err = execute("yum", filepath.Join(path, "wazuh"), "install", "-y", "wazuh-agent")
			if err != nil {
				return err
			}
		}

		if family == "debian" || family == "rhel" {
			err = generateFromTemplate(config, templateFile, configFile)
			if err != nil {
				return err
			}

			err = writeToFile(filepath.Join("/", "var", "ossec", "etc", "client.keys"), string(dKey[:]))
			if err != nil {
				return err
			}

			_, err := execute("systemctl", filepath.Join(path, "beats"), "enable", "wazuh-agent")
			if err != nil {
				return err
			}

			_, err = execute("systemctl", filepath.Join(path, "beats"), "restart", "wazuh-agent")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
