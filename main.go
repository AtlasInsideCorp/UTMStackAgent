package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/AtlasInsideCorp/UTMStackAgent/configuration"
	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
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

var cons = configuration.GetConstConfig()
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
	// Get current path
	path, err := utils.GetMyPath()
	if err != nil {
		fmt.Printf("Failed to get current path: %v", err)
		h.FatalError("Failed to get current path: %v", err)
	}

	// Configuring log saving
	var logger = utils.CreateLogger(filepath.Join(path, "logs", "utmstack_agent.log"))
	defer logger.Close()
	log.SetOutput(logger)

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
