package main

import (
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
				_, err = execute(
					filepath.Join(path, "beats", "windows", "winlogbeat", "winlogbeat.exe"),
					filepath.Join(path, "beats", "windows", "winlogbeat"),
				)
				if err != nil {
					h.FatalError("error running winlogbeat: %v", err)
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
		templateFile := filepath.Join(path, "templates", "winlogbeat.template")
		err := generateFromTemplate(config, templateFile, configFile)
		if err != nil {
			return err
		}
	case "linux":
		configFile := filepath.Join("/", "etc", "filebeat", "filebeat.yml")
		templateFile := filepath.Join(path, "templates", "filebeat-linux.template")

		family, err := detectLinuxFamily()
		if err != nil {
			return err
		}

		switch family {
		case "debian":
			_, err := execute("dpkg", filepath.Join(path, "beats"), "-i", "filebeat-oss-*-amd64.deb")
			if err != nil {
				return err
			}

		case "rhel":
			_, err := execute("yum", filepath.Join(path, "beats"), "install", "-y", "filebeat-oss-*-x86_64.rpm")
			if err != nil {
				return err
			}
		}

		if family == "debian" || family == "rhel" {
			err = generateFromTemplate(config, templateFile, configFile)
			if err != nil {
				return err
			}
			_, err := execute("systemctl", filepath.Join(path, "beats"), "enable", "filebeat")
			if err != nil {
				return err
			}
			_, err = execute("systemctl", filepath.Join(path, "beats"), "restart", "filebeat")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
