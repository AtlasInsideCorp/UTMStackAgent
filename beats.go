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
					h.FatalError("error running beat: %v", err)
				}
			})
		}
	}()
}

func configureBeat(ip string) {
	path, err := getMyPath()
	if err != nil {
		h.FatalError("error getting path: %v", err)
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

	switch runtime.GOOS {
	case "windows":
		configFile := filepath.Join(path, "beats", "windows", "winlogbeat", "winlogbeat.yml")
		templateFile := filepath.Join(path, "templates", "winlogbeat.template")
		config := BeatConfig{ip, ca, clientCert, clientKey}
		err := generateFromTemplate(config, templateFile, configFile)
		h.FatalError("error configuring beat: %v", err)
	}
}
