package main

import (
	"os"
	"sync"
	"time"

	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
)

type config struct {
	Server             string `yaml:"server"`
	AgentID            string `yaml:"agent-id"`
	AgentKey           string `yaml:"agent-key"`
	SkipCertValidation bool   `yaml:"skip-cert-validation"`
}

var oneConfigRead sync.Once
var cnf config

func readConfig() {
	err := utils.ReadYAML("config.yml", &cnf)
	if err != nil {
		h.Error("error reading config %v", err)
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
}

func getConfig() config {
	oneConfigRead.Do(func() { readConfig() })
	return cnf
}
