package utils

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

// loadTLSCredentials loads and returns the TLS transport credentials needed to establish a secure connection to a server.
// Returns an instance of credentials.TransportCredentials that can be used in configuring a gRPC client for secure communication,
// or an error on failure to load or set TLS credentials.
func LoadTLSCredentials(certPath string) (credentials.TransportCredentials, error) {
	// Load the server's certificate
	serverCert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(serverCert) {
		return nil, fmt.Errorf("failed to add server certificate to the certificate pool")
	}

	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}
