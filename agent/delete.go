package agent

import (
	"context"
	"fmt"
	"os"

	"github.com/AtlasInsideCorp/UTMStackAgent/configuration"
	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
	"google.golang.org/grpc/metadata"
)

// DeleteAgent deletes an agent from the gRPC server using the provided configuration information.
func DeleteAgent(cnf configuration.Config, client AgentServiceClient, ctx context.Context) error {
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

	// Delete the agent with the gRPC server
	request := &AgentRequest{
		Ip:       ip,
		Hostname: hostname,
		Os:       osName,
		Platform: osPlatform,
		Version:  osVersion,
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "connection-key", cnf.UTMKey)
	_, err = client.DeleteAgent(ctx, request)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
