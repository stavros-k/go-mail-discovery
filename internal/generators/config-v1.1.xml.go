package generators

import (
	"encoding/xml"

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

type Config_v1_1_xml_params struct {
	Domain              string
	DisplayName         string
	Username            string
	SmtpGlobalPreferred bool
	Provider            providers.Provider
}

func socketType(t providers.SocketType) string {
	switch t {
	case providers.PlainSocketType:
		return "plain"
	case providers.SSLSocketType:
		return "SSL"
	case providers.StartTLSSocketType:
		return "STARTTLS"
	}
	return "SSL"

}

func (c *ClientConfig) Bytes() ([]byte, error) {
	return xml.MarshalIndent(c, "", "\t")
}

func NewConfig_v1_1_xml(p Config_v1_1_xml_params) *ClientConfig {
	incServers := []IncomingServer{}
	if p.Provider.ImapServer != nil {
		incServers = append(incServers, IncomingServer{
			Type:           "imap",
			Username:       p.Username,
			Hostname:       p.Provider.ImapServer.Hostname,
			Port:           p.Provider.ImapServer.Port,
			SocketType:     socketType(p.Provider.ImapServer.SocketType),
			Authentication: p.Provider.ImapServer.Authentication,
		})
	}
	if p.Provider.Pop3Server != nil {
		incServers = append(incServers, IncomingServer{
			Type:           "pop3",
			Username:       p.Username,
			Hostname:       p.Provider.Pop3Server.Hostname,
			Port:           p.Provider.Pop3Server.Port,
			SocketType:     socketType(p.Provider.Pop3Server.SocketType),
			Authentication: p.Provider.Pop3Server.Authentication,
		})
	}

	return &ClientConfig{
		Version: "1.1",
		EmailProvider: EmailProvider{
			ID:              p.Provider.ID,
			Domain:          p.Domain,
			DisplayName:     p.DisplayName,
			IncomingServers: incServers,
			OutgoingServer: OutgoingServer{
				Type:                     "smtp",
				Username:                 p.Username,
				Hostname:                 p.Provider.SmtpServer.Hostname,
				Port:                     p.Provider.SmtpServer.Port,
				SocketType:               socketType(p.Provider.SmtpServer.SocketType),
				Authentication:           p.Provider.SmtpServer.Authentication,
				UseGlobalPreferredServer: p.Provider.SmtpServer.UseGlobalPreferredServer,
			},
		},
	}
}
