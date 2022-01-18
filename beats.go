package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
)

func startBeat() {
	var runOnce sync.Once
	go func() {
		path, err := getMyPath()
		if err != nil {
			h.FatalError("error getting path: %v", err)
		}
		switch runtime.GOOS {
		case "windows":
			runOnce.Do(func() {
				result, err := execute(
					filepath.Join(path, "beats", "windows", "winlogbeat", "winlogbeat.exe"),
					filepath.Join(path, "beats", "windows", "winlogbeat"),
					"--strict.perms=false",
					"-c",
					"winlogbeat.yml",
				)
				if err {
					h.FatalError("error running winlogbeat: %s", result)
				}
			})
		}
	}()
}

func configureBeat(ip string) error {
	path, err := getMyPath()
	if err != nil {
		return err
	}

	clientCert := filepath.Join(path, "keys", TLSCRT)
	clientKey := filepath.Join(path, "keys", TLSKEY)
	ca := filepath.Join(path, "keys", TLSCA)

	type BeatConfig struct {
		IP         string
		CA         string
		ClientCert string
		ClientKey  string
	}

	config := BeatConfig{ip, ca, clientCert, clientKey}

	switch runtime.GOOS {
	case "windows":
		configFile := filepath.Join(path, "beats", "windows", "winlogbeat", "winlogbeat.yml")
		templateFile := filepath.Join(path, "templates", "winlogbeat.yml")
		err := generateFromTemplate(config, templateFile, configFile)
		if err != nil {
			return err
		}
	case "linux":
		configFile := filepath.Join("/", "etc", "filebeat", "filebeat.yml")
		templateFile := filepath.Join(path, "templates", "filebeat-linux.yml")

		family, err := detectLinuxFamily()
		if err != nil {
			return err
		}

		switch family {
		case "debian":
			result, err := execute("dpkg", filepath.Join(path, "beats"), "-i", "filebeat-oss-*-amd64.deb")
			if err {
				return fmt.Errorf("%s", result)
			}
		case "rhel":
			result, err := execute("yum", filepath.Join(path, "beats"), "localinstall", "-y", "filebeat-oss-*-x86_64.rpm")
			if err {
				return fmt.Errorf("%s", result)
			}
		}

		if family == "debian" || family == "rhel" {
			err = generateFromTemplate(config, templateFile, configFile)
			if err != nil {
				return err
			}

			result, err := execute("filebeat", filepath.Join(path, "beats"), "modules", "enable", "system")
			if err {
				return fmt.Errorf("%s", result)
			}

			result, err = execute("systemctl", filepath.Join(path, "beats"), "enable", "filebeat")
			if err {
				return fmt.Errorf("%s", result)
			}

			result, err = execute("systemctl", filepath.Join(path, "beats"), "restart", "filebeat")
			if err {
				return fmt.Errorf("%s", result)
			}
		}
	}
	return nil
}
