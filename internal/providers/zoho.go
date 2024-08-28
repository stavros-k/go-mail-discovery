package providers

func zohoConfig() Provider {
	return Provider{
		ID: "zoho.eu",
		ImapServer: &IncomingServerConfig{
			Hostname:       "imappro.zoho.eu",
			Port:           993,
			SocketType:     SSLSocketType,
			Authentication: CleartextAuthenticationType,
		},
		Pop3Server: &IncomingServerConfig{
			Hostname:       "poppro.zoho.eu",
			Port:           995,
			SocketType:     SSLSocketType,
			Authentication: CleartextAuthenticationType,
		},
		SmtpServer: &OutgoingServerConfig{
			Hostname:                 "smtppro.zoho.eu",
			Port:                     587,
			SocketType:               StartTLSSocketType,
			Authentication:           CleartextAuthenticationType,
			UseGlobalPreferredServer: false,
		},
	}
}
