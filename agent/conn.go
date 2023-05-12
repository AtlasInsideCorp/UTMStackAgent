package agent

import (
	"fmt"
	"path/filepath"

	"github.com/AtlasInsideCorp/UTMStackAgent/configuration"
	"github.com/AtlasInsideCorp/UTMStackAgent/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectToServer(cnf configuration.Config, cons configuration.ConstConfig, serverAddress string) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error
	if cnf.AllowInsecure {
		conn, err = grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(cons.MaxMessageSize)))
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
	} else {
		// Load TLS credentials
		path, err := utils.GetMyPath()
		if err != nil {
			return nil, fmt.Errorf("failed to get current path: %v", err)
		}
		tlsCredentials, err := utils.LoadTLSCredentials(filepath.Join(path, "certs", cons.UTMCRT))
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS credentials: %v", err)
		}
		conn, err = grpc.Dial(serverAddress, grpc.WithTransportCredentials(tlsCredentials), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(cons.MaxMessageSize)))
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
	}

	return conn, nil
}
