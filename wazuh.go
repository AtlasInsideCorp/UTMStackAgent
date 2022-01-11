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
		IP  string
	}

	dKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}

	config := WazuhConfig{ip}

	switch runtime.GOOS {
	case "windows":
		ossecFile := filepath.Join(path, "wazuh", "windows", "ossec.conf")
		ossecTemplateFile := filepath.Join(path, "templates", "wazuh-windows.template")

		err := generateFromTemplate(config, ossecTemplateFile, ossecFile)
		if err != nil {
			return err
		}

		err = writeToFile(filepath.Join(path, "wazuh", "windows", "client.keys"), string(dKey[:]))
		if err != nil {
			return err
		}
	}
	return nil
}
