package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	pb "github.com/AtlasInsideCorp/UTMStackAgent/agent"
	"github.com/AtlasInsideCorp/UTMStackAgent/configuration"
	"github.com/AtlasInsideCorp/UTMStackAgent/stream"
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
			// Read the config from the file
			var cnf configuration.Config
			err = utils.ReadYAML(filepath.Join(path, "config.yml"), &cnf)
			if err != nil {
				fmt.Printf("failed to read config file: %v", err)
				h.FatalError("failed to read config file: %v", err)
			}

			// Connect to the gRPC server
			conn, err := pb.ConnectToServer(cnf, cons, cnf.Server+":"+strconv.Itoa(cons.AGENTMANAGERPORT))
			if err != nil {
				fmt.Printf("Failed to connect to gRPC server: %v", err)
				h.FatalError("Failed to connect to gRPC server: %v", err)
			}
			defer conn.Close()

			// Create a client for AgentService
			agentClient := pb.NewAgentServiceClient(conn)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Start the AgentStream
			stream.StartStream(cnf, agentClient, ctx, cancel, h)
			startBeat()
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
			<-signals

		case "install":
			cnf := configuration.GetInitialConfig()

			// Connect to the gRPC server
			conn, err := pb.ConnectToServer(cnf, cons, cnf.Server+":"+strconv.Itoa(cons.AGENTMANAGERPORT))
			if err != nil {
				fmt.Printf("failed to connect to gRPC server: %v", err)
				h.FatalError("failed to connect to gRPC server: %v", err)
			}
			defer conn.Close()

			// Create a client for AgentService
			agentClient := pb.NewAgentServiceClient(conn)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Register the agent and write the config to the file
			err = pb.RegisterAgent(&cnf, agentClient, ctx)
			if err != nil {
				fmt.Printf("failed to register agent: %v", err)
				h.FatalError("failed to register agent: ", err)
			}

			err = configureBeat(cnf.Server)
			if err != nil {
				h.Error("can't configure beat: %v", err)
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}

			fmt.Println("UMTStack Agent configured correctly")
			h.Info("UMTStack Agent configured correctly")
			os.Exit(0)

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
