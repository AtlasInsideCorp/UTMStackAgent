package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

func registerAgent(endPoint, name string, key string) ([]byte, error) {
	payload := strings.NewReader(fmt.Sprintf("{\n\"agentName\": \"%s\"\n}", name))
	
	req, err := http.NewRequest("POST", endPoint, payload)
	if err != nil {
		return nil, err
	}
	
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("UTM-Token", key)
	
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200{
		return nil, fmt.Errorf("%s", string(body[:]))
	}

	return body, nil
}

type config struct {
	Server   string `yaml:"server"`
	AgentID  string `yaml:"agent-id"`
	AgentKey string `yaml:"agent-key"`
}

var oneConfigRead sync.Once
var cnf config

func readConfig() {
	err := readYAML("config.yml", &cnf)
	if err != nil {
		h.FatalError("error reading config %v", err)
	}
}

func getConfig() config {
	oneConfigRead.Do(func() { readConfig() })
	return cnf
}

func writeConfig(cnf config) error {
	err := writeYAML("config.yml", cnf)
	if err != nil {
		return err
	}
	return nil
}
