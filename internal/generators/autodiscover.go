package generators

import (
	"encoding/xml"
	"fmt"

	"github.com/stavros-k/go-mail-discovery/internal/providers"
)

const AutoDiscoverXmlns = `http://schemas.microsoft.com/exchange/autodiscover/responseschema/2006`
const ResponseXmlns = `http://schemas.microsoft.com/exchange/autodiscover/outlook/responseschema/2006a`

type AutoDiscoverConfig struct {
	XMLName  xml.Name `xml:"Autodiscover"`
	Xmlns    string   `xml:"xmlns,attr"`
	Response Response `xml:"Response"`
}

type Response struct {
	XMLName xml.Name `xml:"Response"`
	Xmlns   string   `xml:"xmlns,attr"`
	User    User     `xml:"User"`
	Account Account  `xml:"Account"`
}

type User struct {
	XMLName     xml.Name `xml:"User"`
	DisplayName string   `xml:"DisplayName"`
}

type Account struct {
	XMLName     xml.Name   `xml:"Account"`
	AccountType string     `xml:"AccountType"`
	Action      string     `xml:"Action"`
	Protocol    []Protocol `xml:"Protocol"`
}

type Protocol struct {
	XMLName            xml.Name          `xml:"Protocol"`
	Type               string            `xml:"Type"`
	Server             string            `xml:"Server"`
	Port               int               `xml:"Port"`
	SecurePasswordAuth *AutoDiscoverBool `xml:"SPA"`
	Encryption         string            `xml:"Encryption"`
	SSL                *AutoDiscoverBool `xml:"SSL,omitempty"`
	AuthRequired       *AutoDiscoverBool `xml:"AuthRequired"`
	UsePOPAuth         *AutoDiscoverBool `xml:"UsePOPAuth,omitempty"`
	SMTPLast           *AutoDiscoverBool `xml:"SMTPLast,omitempty"`
	DomainRequired     *AutoDiscoverBool `xml:"DomainRequired,omitempty"`
	LoginName          string            `xml:"LoginName"`
}

type AutoDiscoverBool bool

func NewAutoDiscoverBoolPtr(b bool) *AutoDiscoverBool {
	result := new(AutoDiscoverBool)
	if b {
		*result = AutoDiscoverBool(true)
	}
	return result

}

func (b AutoDiscoverBool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if b {
		return e.EncodeElement("on", start)
	}
	return e.EncodeElement("off", start)
}

func (a AutoDiscoverConfig) Bytes() ([]byte, error) {
	data, err := xml.MarshalIndent(a, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling plist: %w", err)
	}
	data = append([]byte(xml.Header), data...)
	return data, nil
}

func (a *AutoDiscoverConfig) String() string {
	b, err := a.Bytes()
	if err != nil {
		return fmt.Sprintf("Error generating config: %v", err)
	}
	return string(b)
}

type AutoDiscoverConfigParams struct {
	Domain      string
	DisplayName string
	Username    string
	Provider    providers.Provider
}

func NewAutoDiscoverConfig(p AutoDiscoverConfigParams) (*AutoDiscoverConfig, error) {
	if p.Provider.ID == "" {
		return nil, fmt.Errorf("invalid provider: ID is empty")
	}

	if p.Provider.SmtpServer == nil {
		return nil, fmt.Errorf("no SMTP server configured for provider %s", p.Provider.ID)
	}

	servers := []Protocol{}
	if p.Provider.ImapServer != nil {
		servers = append(servers, createIncomingServerAutoDiscovery("IMAP", p.Username, p.Provider.ImapServer))
	}
	if p.Provider.Pop3Server != nil {
		servers = append(servers, createIncomingServerAutoDiscovery("POP3", p.Username, p.Provider.Pop3Server))
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no incoming servers configured for provider %s", p.Provider.ID)
	}

	servers = append(servers, createOutgoingServerAutoDiscovery(p.Username, p.Provider.SmtpServer))

	return &AutoDiscoverConfig{
		Xmlns: AutoDiscoverXmlns,
		Response: Response{
			Xmlns: ResponseXmlns,
			User: User{
				DisplayName: p.DisplayName,
			},
			Account: Account{
				AccountType: "email",
				Action:      "settings",
				Protocol:    servers,
			},
		},
	}, nil
}

func getEncryption(s providers.SocketType) string {
	switch s {
	case providers.SSLSocketType:
		return "SSL"
	case providers.StartTLSSocketType:
		return "TLS"
	case providers.PlainSocketType:
		return "None"
	default:
		return "Auto"
	}
}

func createIncomingServerAutoDiscovery(serverType string, username string, config *providers.IncomingServerConfig) Protocol {
	return Protocol{
		Type:               serverType,
		Server:             config.Hostname,
		Port:               config.Port,
		SecurePasswordAuth: NewAutoDiscoverBoolPtr(false),
		Encryption:         getEncryption(config.SocketType),
		SSL:                NewAutoDiscoverBoolPtr(config.SocketType == providers.SSLSocketType),
		AuthRequired:       NewAutoDiscoverBoolPtr(config.Authentication != "none"),
		DomainRequired:     NewAutoDiscoverBoolPtr(false),
		LoginName:          username,
	}
}

func createOutgoingServerAutoDiscovery(username string, config *providers.OutgoingServerConfig) Protocol {
	return Protocol{
		Type:               "smtp",
		Server:             config.Hostname,
		Port:               config.Port,
		SecurePasswordAuth: NewAutoDiscoverBoolPtr(false),
		Encryption:         getEncryption(config.SocketType),
		SSL:                NewAutoDiscoverBoolPtr(config.SocketType == providers.SSLSocketType),
		AuthRequired:       NewAutoDiscoverBoolPtr(config.Authentication != "none"),
		UsePOPAuth:         NewAutoDiscoverBoolPtr(false),
		SMTPLast:           NewAutoDiscoverBoolPtr(false),
		DomainRequired:     NewAutoDiscoverBoolPtr(false),
		LoginName:          username,
	}
}
