package configuration

type ConstConfig struct {
	AGENTMANAGERPORT int
	MaxMessageSize   int
	TLSCA            string
	TLSCRT           string
	TLSKEY           string
	UTMCRT           string
	UTMKEY           string
}

// GetConstConfig returns an object with the constant configurations
func GetConstConfig() ConstConfig {
	var cons ConstConfig
	cons.AGENTMANAGERPORT = 50050
	cons.MaxMessageSize = 1024 * 1024 * 1024
	cons.TLSCA = "ca.crt"
	cons.TLSCRT = "client.crt"
	cons.TLSKEY = "client.key"
	cons.UTMCRT = "utm.crt"
	cons.UTMKEY = "utm.key"

	return cons
}
