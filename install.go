package main

import (
	"os"
	"strconv"
	"time"
)

func install(ip, utmKey, skip string) {
	var insecure bool
	if skip == "yes" {
		insecure = true
	}

	hostName, err := os.Hostname()
	if err != nil {
		h.Error("can't get the hostname: %v", err)
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}

	agent, err := registerAgent(AGENTMANAGERPROTO+"://"+ip+":"+strconv.Itoa(AGENTMANAGERPORT), hostName, utmKey, insecure)
	if err != nil {
		h.Error("can't register agent: %v", err)
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}

	cnf := config{Server: ip, AgentID: agent.ID, AgentKey: agent.Key, SkipCertValidation: insecure}
	err = writeConfig(cnf)
	if err != nil {
		h.Error("can't write agent config: %v", err)
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}

	err = configureBeat(ip)
	if err != nil {
		h.Error("can't configure beat: %v", err)
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}

	err = autoStart()
	if err != nil {
		h.Error("can't configure agent service: %v", err)
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}

	os.Exit(0)
}
