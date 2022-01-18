package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var IRMutex sync.Mutex

func incidentResponse() {
	go func() {
		serverName, err := os.Hostname()
		if err != nil {
			h.Error("error getting hostname: %v", err)
		}

		path, err := getMyPath()
		if err != nil {
			h.Error("error getting path: %v", err)
		}

		for {
			IRMutex.Lock()
			cnf := getConfig()
			actions, err := getCommands(
				AGENTMANAGERPROTO+
					"://"+
					cnf.Server+
					":"+
					strconv.Itoa(AGENTMANAGERPORT)+
					GETCOMMANDSENDPOINT+
					fmt.Sprintf("?agentName=%s", serverName),
				cnf.AgentID,
				cnf.AgentKey,
			)
			var commands []struct {
				id      int64
				command string
			}

			if err == nil {
				json.Unmarshal(actions, &commands)
			}

			for _, c := range commands {
				cmd := strings.Split(c.command, " ")
				var response string
				if len(cmd) > 1 {
					response, _ = execute(cmd[0], path, cmd[1:]...)
				} else {
					response, _ = execute(cmd[0], path)
				}

				commandResponse(
					AGENTMANAGERPROTO+
						"://"+
						cnf.Server+
						":"+
						strconv.Itoa(AGENTMANAGERPORT)+
						COMMANDSRESPONSEENDPOINT+
						fmt.Sprintf("?agentName=%s", serverName),
					cnf.AgentID,
					cnf.AgentKey,
					c.id,
					response,
				)
			}

			IRMutex.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()
}

func getCommands(endPoint, agentId, key string) ([]byte, error) {
	req, err := http.NewRequest("GET", endPoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Agent-Id", agentId)
	req.Header.Add("Agent-Key", key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%s", string(body[:]))
	}

	return body, nil
}

func commandResponse(endPoint, agentId, key string, id int64, response string) error {
	payload := strings.NewReader(fmt.Sprintf("{\n\"jobId\": %d,\n\"result\": \"%s\"\n}", id, response))

	req, err := http.NewRequest("POST", endPoint, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Agent-Id", agentId)
	req.Header.Add("Agent-Key", key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("%s", string(body[:]))
	}

	return nil
}
