package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quantfall/holmes"
)

const (
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

type jobResult struct {
	JobId  int64  `json:"jobId"`
	Result string `json:"result"`
}

func main() {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "run":
			incidentResponse()
			startBeat()
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
			<-signals

		case "install":
			var ip string
			var utmKey string
			var skip string

			fmt.Println("Manager IP or FQDN:")
			if _, err := fmt.Scanln(&ip); err != nil {
				h.Error("can't get the manager IP or FQDN: %v", err)
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}

			fmt.Println("Registration Key:")
			if _, err := fmt.Scanln(&utmKey); err != nil {
				h.Error("can't get the registration key: %v", err)
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}

			fmt.Println("Skip certificate validation (yes or no):")
			if _, err := fmt.Scanln(&skip); err != nil {
				h.Error("can't get certificate validation response: %v", err)
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}

			install(ip, utmKey, skip)

		case "silent-install":
			ip := os.Args[2]
			utmKey := os.Args[3]
			skip := os.Args[4]

			install(ip, utmKey, skip)

		default:
			fmt.Println("unknown option")
		}
	} else {
		err := uninstall()
		if err != nil {
			h.Error("can't remove agent dependencies or configurations: %v", err)
			time.Sleep(10 * time.Second)
			os.Exit(1)
		}

		os.Exit(0)
	}
}
