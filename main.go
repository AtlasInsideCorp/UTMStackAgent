package main

import (
	"crypto/tls"
	"encoding/json"
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
	TLSSKIPVERIFICATION      = false
	AGENTMANAGERPROTO        = "http"
	AGENTMANAGERPORT         = 8080
	REGISTRATIONENDPOINT     = "/api/v1/agent"
	GETCOMMANDSENDPOINT      = "/api/v1/agent-by-name"
	COMMANDSRESPONSEENDPOINT = "/api/v1/agent-by-name"
	TLSCA                    = "ca.crt"
	TLSCRT                   = "client.crt"
	TLSKEY                   = "client.key"
)

var h = holmes.New("debug", "UTMStack")

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: TLSSKIPVERIFICATION}
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
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
		fmt.Println("Please insert the Master or Proxy IP")
		if _, err := fmt.Scanln(&ip); err != nil {
			h.FatalError("can't get Master or Proxy ip addr: %v", err)
		}
		if net.ParseIP(ip) == nil {
			h.FatalError("%v is not a valid IP", ip)
		}
		fmt.Println("Please insert the UTMStack Key")
		if _, err := fmt.Scanln(&utmKey); err != nil {
			h.FatalError("can't get the UTMStack Key: %v", err)
		}
		hostName, err := os.Hostname()
		if err != nil {
			h.FatalError("can't get the hostname: %v", err)
		}
		regReq, err := registerAgent(
			AGENTMANAGERPROTO+
				"://"+
				ip+
				":"+
				strconv.Itoa(AGENTMANAGERPORT)+
				REGISTRATIONENDPOINT,
			hostName,
			utmKey,
		)
		if err != nil {
			h.FatalError("can't register agent: %v", err)
		}
		var agentDetails struct {
			ID  string `json:"id"`
			Key string `json:"key"`
		}
		err = json.Unmarshal(regReq, &agentDetails)
		if err != nil {
			h.FatalError("can't decode agent details: %v", err)
		}
		cnf := config{Server: ip, AgentID: agentDetails.ID, AgentKey: agentDetails.Key}
		writeConfig(cnf)
		configureBeat(ip)
		configureWazuh(ip, cnf.AgentKey)
		err = autoStart()
		if err != nil {
			h.FatalError("can't configure agent service: %v", err)
		}
	}
}
