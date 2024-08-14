package providers

func Zoho() Server {
	return Server{
		ID: "zoho.eu",
		ImapServer: &incomingServerConfig{
			Hostname:       "imappro.zoho.eu",
			Port:           993,
			SocketType:     "ssl",
			Authentication: "password-cleartext",
		},
		Pop3Server: &incomingServerConfig{
			Hostname:       "poppro.zoho.eu",
			Port:           995,
			SocketType:     "ssl",
			Authentication: "password-cleartext",
		},
		SmtpServer: &outgoingServerConfig{
			Hostname:                 "smtppro.zoho.eu",
			Port:                     587,
			SocketType:               "starttls",
			Authentication:           "password-cleartext",
			UseGlobalPreferredServer: false,
		},
	}
}
