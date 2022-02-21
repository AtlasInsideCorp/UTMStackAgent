package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func startWazuh() {
	var runOnce sync.Once
	go func() {
		path, err := getMyPath()
		if err != nil {
			h.Error("error getting path: %v", err)
			time.Sleep(10 * time.Second)
			os.Exit(1)
		}
		switch runtime.GOOS {
		case "windows":
			runOnce.Do(func() {
				execute(
					filepath.Join(path, "wazuh", "windows", "wazuh-agent.exe"),
					filepath.Join(path, "wazuh", "windows"),
					"install-service",
				)
				result, errB := execute(
					filepath.Join(path, "nssm.exe"),
					path,
					"start",
					"WazuhSvc",
				)
				if errB {
					h.Error("error running wazuh: %s", result)
					time.Sleep(10 * time.Second)
					os.Exit(1)
				}
			})
		}
	}()
}

func stopWazuh() {
	var runOnce sync.Once
	path, err := getMyPath()
	if err != nil {
		h.Error("error getting path: %v", err)
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	switch runtime.GOOS {
	case "windows":
		runOnce.Do(func() {
			result, errB := execute(
				filepath.Join(path, "nssm.exe"),
				path,
				"stop",
				"WazuhSvc",
			)
			if errB {
				h.Error("error stopping wazuh: %s", result)
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}
			execute(
				filepath.Join(path, "wazuh", "windows", "wazuh-agent.exe"),
				filepath.Join(path, "wazuh", "windows"),
				"uninstall-service",
			)
		})
	}
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

			result, errB := execute("apt", path, "update")
			if errB {
				return fmt.Errorf("%s", result)
			}

			result, errB = execute("apt", path, "install", "-y", "curl", "apt-transport-https", "lsb-release", "gnupg2", "wget")
			if errB {
				return fmt.Errorf("%s", result)
			}

			err := download("https://packages.wazuh.com/key/GPG-KEY-WAZUH")
			if err != nil {
				return err
			}

			result, errB = execute("apt-key", path, "add", "GPG-KEY-WAZUH")
			if errB {
				return fmt.Errorf("%s", result)
			}

			err = writeToFile(filepath.Join("/", "etc", "apt", "sources.list.d", "wazuh.list"), "deb https://packages.wazuh.com/4.x/apt/ stable main")
			if err != nil {
				return err
			}

			result, errB = execute("apt", path, "update")
			if errB {
				return fmt.Errorf("%s", result)
			}

			result, errB = execute("apt", path, "install", "-y", "wazuh-agent")
			if errB {
				return fmt.Errorf("%s", result)
			}

		case "rhel":
			templateFile = filepath.Join(path, "templates", "wazuh-rhel.conf")

			result, errB := execute("rpm", path, "--import", "https://packages.wazuh.com/key/GPG-KEY-WAZUH")
			if errB {
				return fmt.Errorf("%s", result)
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

			result, errB = execute("yum", path, "install", "-y", "wazuh-agent")
			if errB {
				return fmt.Errorf("%s", result)
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

			result, errB := execute("systemctl", path, "enable", "wazuh-agent")
			if errB {
				return fmt.Errorf("%s", result)
			}

			result, errB = execute("systemctl", path, "restart", "wazuh-agent")
			if errB {
				return fmt.Errorf("%s", result)
			}
		}
	}
	return nil
}
