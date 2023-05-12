package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
)

func startBeat() {
	var runOnce sync.Once
	go func() {
		path, err := utils.GetMyPath()
		if err != nil {
			h.Error("error getting path: %v", err)
			time.Sleep(10 * time.Second)
			os.Exit(1)
		}
		switch runtime.GOOS {
		case "windows":
			runOnce.Do(func() {
				result, err := utils.Execute(
					filepath.Join(path, "beats", "winlogbeat", "winlogbeat.exe"),
					filepath.Join(path, "beats", "winlogbeat"),
					"--strict.perms=false",
					"-c",
					"winlogbeat.yml",
				)
				if err {
					h.Error("error running winlogbeat: %s", result)
					time.Sleep(10 * time.Second)
					os.Exit(1)
				}
			})
		}
	}()
}

func configureBeat(ip string) error {
	path, err := utils.GetMyPath()
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
		configFile := filepath.Join(path, "beats", "winlogbeat", "winlogbeat.yml")
		templateFile := filepath.Join(path, "templates", "winlogbeat.yml")
		err := utils.GenerateFromTemplate(config, templateFile, configFile)
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
			result, err := utils.Execute("dpkg", filepath.Join(path, "beats"), "-i", filepath.Join(path, "beats", "filebeat-oss-7.13.4-amd64.deb"))
			if err {
				return fmt.Errorf("%s", result)
			}
		case "rhel":
			result, err := utils.Execute("yum", filepath.Join(path, "beats"), "localinstall", "-y", filepath.Join(path, "beats", "filebeat-oss-7.13.4-x86_64.rpm"))
			if err {
				return fmt.Errorf("%s", result)
			}
		}

		if family == "debian" || family == "rhel" {
			err = utils.GenerateFromTemplate(config, templateFile, configFile)
			if err != nil {
				return err
			}

			result, err := utils.Execute("filebeat", filepath.Join(path, "beats"), "modules", "enable", "system")
			if err {
				return fmt.Errorf("%s", result)
			}

			result, err = utils.Execute("systemctl", filepath.Join(path, "beats"), "enable", "filebeat")
			if err {
				return fmt.Errorf("%s", result)
			}

			result, err = utils.Execute("systemctl", filepath.Join(path, "beats"), "restart", "filebeat")
			if err {
				return fmt.Errorf("%s", result)
			}
		}
	}
	return nil
}
