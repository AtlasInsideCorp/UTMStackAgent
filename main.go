package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/quantfall/holmes"
)

const (
	TLSSKIPVERIFICATION      = true
	AGENTMANAGERPROTO        = "https"
	AGENTMANAGERPORT         = 9000
	REGISTRATIONENDPOINT     = "/api/v1/agent"
	GETIDANDKEYENDPOINT      = "/api/v1/agent-id-key-by-name"
	GETCOMMANDSENDPOINT      = "/api/v1/incident-commands"
	COMMANDSRESPONSEENDPOINT = "/api/v1/incident-command/result"
	TLSCA                    = "ca.crt"
	TLSCRT                   = "client.crt"
	TLSKEY                   = "client.key"
)

var h = holmes.New("debug", "UTMStack")

type agentDetails struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: TLSSKIPVERIFICATION}
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "uninstall":
			err := uninstall()
			if err != nil {
				h.FatalError("can't remove agent dependencies or configurations: %v", err)
			}
			
			os.Exit(0)
		case "run":
			incidentResponse()
			startBeat()
			startWazuh()
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
			<-signals
		default:
			fmt.Println("unknown option")
		}
	} else {
		var ip string
		var utmKey string

		fmt.Println("Manager IP or FQDN:")
		if _, err := fmt.Scanln(&ip); err != nil {
			h.FatalError("can't get the manager IP or FQDN: %v", err)
		}

		fmt.Println("Registration Key:")
		if _, err := fmt.Scanln(&utmKey); err != nil {
			h.FatalError("can't get the registration key: %v", err)
		}

		hostName, err := os.Hostname()
		if err != nil {
			h.FatalError("can't get the hostname: %v", err)
		}

		agent, err := registerAgent(AGENTMANAGERPROTO+"://"+ip+":"+strconv.Itoa(AGENTMANAGERPORT), hostName, utmKey,
		)
		if err != nil {
			h.FatalError("Can't register agent: %v", err)
		}

		cnf := config{Server: ip, AgentID: agent.ID, AgentKey: agent.Key}
		err = writeConfig(cnf)
		if err != nil {
			h.FatalError("can't write agent config: %v", err)
		}

		err = configureBeat(ip)
		if err != nil {
			h.FatalError("can't configure beat: %v", err)
		}

		err = configureWazuh(ip, cnf.AgentKey)
		if err != nil {
			h.FatalError("can't configure wazuh: %v", err)
		}

		err = autoStart()
		if err != nil {
			h.FatalError("can't configure agent service: %v", err)
		}

		os.Exit(0)
	}
}
