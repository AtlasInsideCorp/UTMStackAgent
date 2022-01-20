package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

func registerAgent(endPoint, name string, key string) (agentDetails, error) {
	var body []byte
	payload := strings.NewReader(fmt.Sprintf(`{"agentName": "%s"}`, name))

	regReq, err := http.NewRequest("POST", endPoint+REGISTRATIONENDPOINT, payload)
	if err != nil {
		return agentDetails{}, err
	}

	regReq.Header.Add("Content-Type", "application/json")
	regReq.Header.Add("UTM-Token", key)

	regRes, err := http.DefaultClient.Do(regReq)
	if err != nil {
		return agentDetails{}, err
	}
	defer regRes.Body.Close()

	body, err = ioutil.ReadAll(regRes.Body)
	if err != nil {
		return agentDetails{}, err
	}

	if regRes.StatusCode != 200 {
		keyReq, err := http.NewRequest("GET", endPoint+GETIDANDKEYENDPOINT+fmt.Sprintf("?agentName=%s", name), nil)
		if err != nil {
			return agentDetails{}, err
		}

		keyReq.Header.Add("UTM-Token", key)

		keyRes, err := http.DefaultClient.Do(keyReq)
		if err != nil {
			return agentDetails{}, err
		}
		defer keyRes.Body.Close()

		body, err = ioutil.ReadAll(keyRes.Body)
		if err != nil {
			return agentDetails{}, err
		}

		var agentList []agentDetails

		err = json.Unmarshal(body, &agentList)
		if err != nil {
			h.FatalError("can't decode agent details: %v", err)
		}
		h.Debug("Agent Details: %v", agentList)

		return agentList[0], nil
	}

	var agent agentDetails

	err = json.Unmarshal(body, &agent)
	if err != nil {
		h.FatalError("can't decode agent details: %v", err)
	}
	h.Debug("Agent Details: %v", agent)

	return agent, nil
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
