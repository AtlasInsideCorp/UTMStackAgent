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

	pb "github.com/AtlasInsideCorp/UTMStackAgent/agent"
	"github.com/AtlasInsideCorp/UTMStackAgent/beat"
	"github.com/AtlasInsideCorp/UTMStackAgent/configuration"
	"github.com/AtlasInsideCorp/UTMStackAgent/stream"
	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
	"github.com/quantfall/holmes"
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

			// Install and configure the necessary Beats
			beat.InstallBeats(cnf.Server, cons, h)

			fmt.Println("UMTStack Agent configured correctly")
			h.Info("UMTStack Agent configured correctly")
			os.Exit(0)

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

			// Runs the necessary Beats
			beat.RunBeats(h)

			// Start the AgentStream
			stream.StartStream(cnf, agentClient, ctx, cancel, h)
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
			<-signals

		case "uninstall-dependencies":
			// Disable installed beats
			beat.DisableBeats(h)
			os.Exit(0)

		default:
			fmt.Println("unknown option")
			h.Info("unknown option")
		}
	} else {
		// Read the config from the file
		//var cnf configuration.Config
		//err = utils.ReadYAML(filepath.Join(path, "config.yml"), &cnf)
		//if err != nil {
		//	fmt.Printf("failed to read config file: %v", err)
		//	h.FatalError("failed to read config file: %v", err)
		//}

		// Connect to the gRPC server
		//conn, err := configuration.ConnectToServer(cnf, cons, cnf.Server+":"+strconv.Itoa(cons.AGENTMANAGERPORT))
		//if err != nil {
		//	fmt.Printf("Failed to connect to gRPC server: %v", err)
		//	h.FatalError("failed to connect to gRPC server: %v", err)
		//}
		//defer conn.Close()

		// Create a client for AgentService
		//agentClient := pb.NewAgentServiceClient(conn)
		//ctx, cancel := context.WithCancel(context.Background())
		//defer cancel()

		// Delete Agent in the gRPC Server
		//err = configuration.DeleteAgent(cnf, agentClient, ctx)
		//if err != nil {
		//	fmt.Printf("failed to delete agent: %v", err)
		//	h.FatalError("failed to delete agent: %v", err)
		//}

		// Disable installed beats
		beat.DisableBeats(h)
		os.Remove(filepath.Join(path, "config.yml"))
		os.Exit(0)
	}
}
