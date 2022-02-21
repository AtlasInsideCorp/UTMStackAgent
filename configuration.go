package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func registerAgent(endPoint, name string, key string, insecure bool) (agentDetails, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: insecure}
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
			h.Error("can't decode agent details: %v", err)
			time.Sleep(10 * time.Second)
			os.Exit(1)
		}
		h.Debug("Agent Details: %v", agentList)

		return agentList[0], nil
	}

	var agent agentDetails

	err = json.Unmarshal(body, &agent)
	if err != nil {
		h.Error("can't decode agent details: %v", err)
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	h.Debug("Agent Details: %v", agent)

	return agent, nil
}

type config struct {
	Server             string `yaml:"server"`
	AgentID            string `yaml:"agent-id"`
	AgentKey           string `yaml:"agent-key"`
	SkipCertValidation bool   `yaml:"skip-cert-validation"`
}

var oneConfigRead sync.Once
var cnf config

func readConfig() {
	err := readYAML("config.yml", &cnf)
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

func writeConfig(cnf config) error {
	err := writeYAML("config.yml", cnf)
	if err != nil {
		return err
	}
	return nil
}
