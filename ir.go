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
				response, _ := execute(c.command, path)
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
			time.Sleep(10 * time.Second)
		}
	}()
}

func getCommands(endPoint, agentId, key string) ([]byte, error) {
	var err error
	if req, err := http.NewRequest("GET", endPoint, nil); err == nil {
		req.Header.Add("Agent-Id", agentId)
		req.Header.Add("Agent-Key", key)
		if res, err := http.DefaultClient.Do(req); err == nil {
			defer res.Body.Close()
			if body, err := ioutil.ReadAll(res.Body); err == nil {
				return body, nil
			}
		}
	}
	return nil, err
}

func commandResponse(endPoint, agentId, key string, id int64, response string) error {
	var err error
	payload := strings.NewReader(fmt.Sprintf("{\n\"jobId\": %d,\n\"result\": \"%s\"\n}", id, response))
	if req, err := http.NewRequest("POST", endPoint, payload); err == nil {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Agent-Id", agentId)
		req.Header.Add("Agent-Key", key)
		if res, err := http.DefaultClient.Do(req); err == nil {
			defer res.Body.Close()
			if _, err := ioutil.ReadAll(res.Body); err == nil {
				return nil
			}
		}
	}
	return err
}
