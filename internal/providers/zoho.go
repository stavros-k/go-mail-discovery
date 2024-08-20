package providers

func zohoConfig() Provider {
	return Provider{
		ID: "zoho.eu",
		ImapServer: &IncomingServerConfig{
			Hostname:       "imappro.zoho.eu",
			Port:           993,
			SocketType:     SSLSocketType,
			Authentication: "password-cleartext",
		},
		Pop3Server: &IncomingServerConfig{
			Hostname:       "poppro.zoho.eu",
			Port:           995,
			SocketType:     SSLSocketType,
			Authentication: "password-cleartext",
		},
		SmtpServer: &OutgoingServerConfig{
			Hostname:                 "smtppro.zoho.eu",
			Port:                     587,
			SocketType:               StartTLSSocketType,
			Authentication:           "password-cleartext",
			UseGlobalPreferredServer: false,
		},
	}
}

// For future dynamic configuration:
// func LoadProvidersFromConfig(configPath string) error {
//     // Implementation to load providers from YAML/JSON
// }
