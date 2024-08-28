package providers

import (
	"fmt"
	"strings"
)

type AuthenticationType int

const (
	NoneAuthenticationType AuthenticationType = iota
	CleartextAuthenticationType
)

var authenticationTypes = map[string]AuthenticationType{
	"none":               NoneAuthenticationType,
	"password-cleartext": CleartextAuthenticationType,
}

func (at AuthenticationType) String() string {
	for k, v := range authenticationTypes {
		if v == at {
			return k
		}
	}
	return ""
}

func (at *AuthenticationType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	parsed, ok := authenticationTypes[s]
	if !ok {
		return fmt.Errorf("invalid authentication type: %s", s)
	}
	*at = parsed
	return nil
}

func (at *AuthenticationType) IsValid() bool {
	_, ok := authenticationTypes[at.String()]
	return ok
}

type SocketType int

const (
	PlainSocketType SocketType = iota
	SSLSocketType
	StartTLSSocketType
)

var socketTypes = map[string]SocketType{
	"PLAIN":    PlainSocketType,
	"SSL":      SSLSocketType,
	"STARTTLS": StartTLSSocketType,
}

// String implements the Stringer interface for SocketType
func (st SocketType) String() string {
	for k, v := range socketTypes {
		if v == st {
			return k
		}
	}
	return ""
}

func (st *SocketType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	parsed, ok := socketTypes[s]
	if !ok {
		return fmt.Errorf("invalid socket type: %s", s)
	}
	*st = parsed
	return nil
}
func (st *SocketType) IsValid() bool {
	_, ok := socketTypes[st.String()]
	return ok
}

type Providers struct {
	Providers []Provider `yaml:"providers"`
}

type Provider struct {
	ID         string                `yaml:"id"`
	ImapServer *IncomingServerConfig `yaml:"imap_server"`
	Pop3Server *IncomingServerConfig `yaml:"pop3_server"`
	SmtpServer *OutgoingServerConfig `yaml:"smtp_server"`
}

func (p *Provider) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("provider ID is required")
	}
	if p.ImapServer != nil {
		if err := p.ImapServer.Validate(); err != nil {
			return fmt.Errorf("imap server validation failed for provider %s: %v", p.ID, err)
		}
	}
	if p.Pop3Server != nil {
		if err := p.Pop3Server.Validate(); err != nil {
			return fmt.Errorf("pop3 server validation failed for provider %s: %v", p.ID, err)
		}
	}
	if p.SmtpServer != nil {
		if err := p.SmtpServer.Validate(); err != nil {
			return fmt.Errorf("smtp server validation failed for provider %s: %v", p.ID, err)
		}
	}
	return nil
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

type IncomingServerConfig struct {
	Hostname       string             `yaml:"hostname"`
	Port           int                `yaml:"port"`
	SocketType     SocketType         `yaml:"socket_type"`
	Authentication AuthenticationType `yaml:"authentication"`
}

func (i *IncomingServerConfig) Validate() error {
	if i.Hostname == "" {
		return fmt.Errorf("IncomingServerConfig: hostname is required")
	}
	if i.Port == 0 {
		return fmt.Errorf("IncomingServerConfig: port is required")
	}
	if !i.SocketType.IsValid() {
		return fmt.Errorf("IncomingServerConfig: socket_type [%s] is invalid", i.SocketType)
	}
	if !i.Authentication.IsValid() {
		return fmt.Errorf("IncomingServerConfig: authentication [%s] is invalid", i.Authentication)
	}
	return nil
}

type OutgoingServerConfig struct {
	Hostname                 string             `yaml:"hostname"`
	Port                     int                `yaml:"port"`
	SocketType               SocketType         `yaml:"socket_type"`
	Authentication           AuthenticationType `yaml:"authentication"`
	UseGlobalPreferredServer bool               `yaml:"use_global_preferred_server"`
}

func (o *OutgoingServerConfig) Validate() error {
	if o.Hostname == "" {
		return fmt.Errorf("OutgoingServerConfig: hostname is required")
	}
	if o.Port == 0 {
		return fmt.Errorf("OutgoingServerConfig: port is required")
	}
	if !o.SocketType.IsValid() {
		return fmt.Errorf("OutgoingServerConfig: socket_type [%s] is invalid", o.SocketType)
	}
	if !o.Authentication.IsValid() {
		return fmt.Errorf("OutgoingServerConfig: authentication [%s] is invalid", o.Authentication)
	}
	return nil
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
