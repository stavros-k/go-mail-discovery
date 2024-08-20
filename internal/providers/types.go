package providers

import (
	"fmt"
	"strings"
)

type SocketType int

const (
	PlainSocketType SocketType = iota
	SSLSocketType
	StartTLSSocketType
)

// String implements the Stringer interface for SocketType
func (st SocketType) String() string {
	switch st {
	case PlainSocketType:
		return "Plain"
	case SSLSocketType:
		return "SSL"
	case StartTLSSocketType:
		return "StartTLS"
	default:
		return fmt.Sprintf("Unknown SocketType(%d)", int(st))
	}
}

type Provider struct {
	ID         string
	ImapServer *IncomingServerConfig
	Pop3Server *IncomingServerConfig
	SmtpServer *OutgoingServerConfig
}

type IncomingServerConfig struct {
	Hostname       string
	Port           int
	SocketType     SocketType
	Authentication string
}

type OutgoingServerConfig struct {
	Hostname                 string
	Port                     int
	SocketType               SocketType
	Authentication           string
	UseGlobalPreferredServer bool
}

func (o OutgoingServerConfig) String() string {
	return fmt.Sprintf("Outgoing Server: %s:%d | (%s) | (%s)",
		o.Hostname, o.Port, o.SocketType, o.Authentication,
	)
}

func (i IncomingServerConfig) String() string {
	return fmt.Sprintf("Incoming Server: %s:%d | (%s) | (%s)",
		i.Hostname, i.Port, i.SocketType, i.Authentication,
	)
}

func (p Provider) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("Provider: %s\n", p.ID))
	if p.ImapServer != nil {
		s.WriteString(fmt.Sprintf("Incoming Server: %s\n", p.ImapServer))
	}
	if p.Pop3Server != nil {
		s.WriteString(fmt.Sprintf("Incoming Server: %s\n", p.Pop3Server))
	}
	if p.SmtpServer != nil {
		s.WriteString(fmt.Sprintf("Outgoing Server: %s\n", p.SmtpServer))
	}

	return s.String()
}
