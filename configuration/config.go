package configuration

import (
	"os"

	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
)

type Config struct {
	Server        string `yaml:"server"`
	UTMKey        string `yaml:"utm-key"`
	AgentID       uint   `yaml:"agent-id"`
	AgentToken    string `yaml:"agent-token"`
	AllowInsecure bool   `yaml:"insecure"`
}

var ConfInstance Config

// GetInitialConfig returns a configuration instance.
func GetInitialConfig() Config {
	ConfInstance.Server = os.Args[2]
	ConfInstance.UTMKey = os.Args[3]
	if os.Args[4] == "yes" {
		ConfInstance.AllowInsecure = true
	}

	return ConfInstance
}

// WriteConfig writes the configuration to a config.yml file
func WriteConfig(cnf *Config) error {
	err := utils.WriteYAML("config.yml", &cnf)
	if err != nil {
		return err
	}
	return nil
}
