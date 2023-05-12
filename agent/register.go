package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AtlasInsideCorp/UTMStackAgent/configuration"
	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
	"google.golang.org/grpc/metadata"
)

// RegisterAgent registers an agent to the gRPC server using the provided configuration information.
func RegisterAgent(cnf *configuration.Config, client AgentServiceClient, ctx context.Context) error {
	path, err := utils.GetMyPath()
	if err != nil {
		return fmt.Errorf("error getting current path: %v", err)
	}

	// Check if the config.yml file exists
	configFilePath := filepath.Join(path, "config.yml")
	if _, err := os.Stat(configFilePath); err == nil {
		// Read the config from the file
		err := utils.ReadYAML(configFilePath, &cnf)
		if err != nil {
			return fmt.Errorf("failed to read config file: %v", err)
		}
	} else {
		// Get the agent config information
		ip, err := utils.GetIPAddress()
		if err != nil {
			return fmt.Errorf("failed to get IP address: %v", err)
		}

		hostname, err := os.Hostname()
		if err != nil {
			return fmt.Errorf("failed to get hostname: %v", err)
		}

		osName, osVersion, osPlatform := utils.GetOSInfo()

		// Register the agent with the gRPC server
		request := &AgentRequest{
			Ip:       ip,
			Hostname: hostname,
			Os:       osName,
			Platform: osPlatform,
			Version:  osVersion,
		}

		ctx = metadata.AppendToOutgoingContext(ctx, "connection-key", cnf.UTMKey)
		response, err := client.RegisterAgent(ctx, request)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		cnf.AgentID = uint(response.Id)
		cnf.AgentToken = response.Token

		// Write the config to the file
		err = configuration.WriteConfig(cnf)
		if err != nil {
			return fmt.Errorf("can't write agent config: %v", err)
		}
	}

	return nil
}
