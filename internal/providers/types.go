package providers

type Server struct {
	ID         string
	ImapServer *incomingServerConfig
	Pop3Server *incomingServerConfig
	SmtpServer *outgoingServerConfig
}

type incomingServerConfig struct {
	Hostname       string
	Port           int
	SocketType     string
	Authentication string
}

type outgoingServerConfig struct {
	Hostname                 string
	Port                     int
	SocketType               string
	Authentication           string
	UseGlobalPreferredServer bool
}
