package generators

import (
	"encoding/xml"
	"fmt"

	"github.com/stavros-k/go-mail-discovery/internal/providers"
)

type ClientConfig struct {
	XMLName       xml.Name      `xml:"clientConfig"`
	Version       string        `xml:"version,attr"`
	EmailProvider EmailProvider `xml:"emailProvider"`
}

type EmailProvider struct {
	ID              string           `xml:"id,attr"`
	Domain          string           `xml:"domain"`
	DisplayName     string           `xml:"displayName"`
	IncomingServers []IncomingServer `xml:"incomingServer"`
	OutgoingServer  OutgoingServer   `xml:"outgoingServer"`
}

type IncomingServer struct {
	Type           string `xml:"type,attr"`
	Hostname       string `xml:"hostname"`
	Port           int    `xml:"port"`
	SocketType     string `xml:"socketType"`
	Authentication string `xml:"authentication"`
	Username       string `xml:"username"`
}

type OutgoingServer struct {
	Type                     string `xml:"type,attr"`
	Hostname                 string `xml:"hostname"`
	Port                     int    `xml:"port"`
	SocketType               string `xml:"socketType"`
	Username                 string `xml:"username"`
	Authentication           string `xml:"authentication"`
	UseGlobalPreferredServer bool   `xml:"useGlobalPreferredServer"`
}

type ConfigV1_1Params struct {
	Domain      string
	DisplayName string
	Username    string
	Provider    providers.Provider
}

func (c *ClientConfig) Bytes() ([]byte, error) {
	return xml.MarshalIndent(c, "", "  ")
}

func NewConfigV1_1(p ConfigV1_1Params) (*ClientConfig, error) {
	if p.Provider.ID == "" {
		return nil, fmt.Errorf("invalid provider: ID is empty")
	}

	incServers := []IncomingServer{}
	if p.Provider.ImapServer != nil {
		incServers = append(incServers, createIncomingServer("imap", p.Username, p.Provider.ImapServer))
	}
	if p.Provider.Pop3Server != nil {
		incServers = append(incServers, createIncomingServer("pop3", p.Username, p.Provider.Pop3Server))
	}

	if len(incServers) == 0 {
		return nil, fmt.Errorf("no incoming servers configured for provider %s", p.Provider.ID)
	}

	if p.Provider.SmtpServer == nil {
		return nil, fmt.Errorf("no SMTP server configured for provider %s", p.Provider.ID)
	}

	return &ClientConfig{
		Version: "1.1",
		EmailProvider: EmailProvider{
			ID:              p.Provider.ID,
			Domain:          p.Domain,
			DisplayName:     p.DisplayName,
			IncomingServers: incServers,
			OutgoingServer:  createOutgoingServer(p.Username, p.Provider.SmtpServer),
		},
	}, nil
}

func createIncomingServer(serverType string, username string, config *providers.IncomingServerConfig) IncomingServer {
	return IncomingServer{
		Type:           serverType,
		Username:       username,
		Hostname:       config.Hostname,
		Port:           config.Port,
		SocketType:     config.SocketType.String(),
		Authentication: config.Authentication,
	}
}

func createOutgoingServer(username string, config *providers.OutgoingServerConfig) OutgoingServer {
	return OutgoingServer{
		Type:                     "smtp",
		Username:                 username,
		Hostname:                 config.Hostname,
		Port:                     config.Port,
		SocketType:               config.SocketType.String(),
		Authentication:           config.Authentication,
		UseGlobalPreferredServer: config.UseGlobalPreferredServer,
	}
}
