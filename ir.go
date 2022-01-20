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
				ID      int64  `json:"id"`
				Command string `json:"command"`
			}

			if err == nil {
				json.Unmarshal(actions, &commands)
			} else {
				h.Error("Error getting commands: %v", err)
			}

			for _, c := range commands {
				cmd := strings.Split(c.Command, " ")
				var response string
				h.Debug("Executing command: %v", cmd)
				if len(cmd) > 1 {
					response, _ = execute(cmd[0], path, cmd[1:]...)
				} else {
					response, _ = execute(cmd[0], path)
				}
				err := commandResponse(
					AGENTMANAGERPROTO+
						"://"+
						cnf.Server+
						":"+
						strconv.Itoa(AGENTMANAGERPORT)+
						COMMANDSRESPONSEENDPOINT+
						fmt.Sprintf("?agentName=%s", serverName),
					cnf.AgentID,
					cnf.AgentKey,
					c.ID,
					response,
				)
				if err != nil {
					h.Error("Error sending command response: %v", err)
				}
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
	result, err := json.Marshal(jobResult{JobId:id, Result:response})
	if err != nil{
		return err
	}
	payload := strings.NewReader(string(result))
	h.Debug("Command execution result: %s", string(result))
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
