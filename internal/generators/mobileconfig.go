package generators

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/google/uuid"
	"github.com/stavros-k/go-mail-discovery/internal/providers"
)

const plistHeader = `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">` + "\n"
const plistVersion = "1.0"
const payloadType = "com.apple.mail.managed"

type MobileConfig struct {
	XMLName xml.Name `xml:"plist"`
	Version string   `xml:"version,attr"`
	Dict    Dict     `xml:"dict"`
}

func (m *MobileConfig) Bytes() ([]byte, error) {
	data, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling plist: %w", err)
	}
	data = bytes.ReplaceAll(data, []byte("></true>"), []byte(" />"))
	data = bytes.ReplaceAll(data, []byte("></false>"), []byte(" />"))
	result := []byte(xml.Header)
	result = append(result, []byte(plistHeader)...)
	result = append(result, data...)
	return result, nil
}

func (m *MobileConfig) String() string {
	b, err := m.Bytes()
	if err != nil {
		return fmt.Sprintf("Error generating config: %v", err)
	}
	return string(b)
}

type Dict struct {
	XMLName xml.Name   `xml:"dict"`
	Key     string     `xml:"key,omitempty"`
	Array   *Array     `xml:"array,omitempty"`
	Items   []DictItem `xml:",any"`
}

type DictItem interface {
	isDictItem()
}

func (a Array) isDictItem()        {}
func (s StringValue) isDictItem()  {}
func (i IntegerValue) isDictItem() {}
func (b BooleanValue) isDictItem() {}

type Array struct {
	XMLName xml.Name `xml:"array"`
	Items   []Dict   `xml:"dict"`
}

type StringValue struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func NewStringEntry(key, value string) []DictItem {
	result := []DictItem{
		StringValue{XMLName: xml.Name{Local: "key"}, Value: key},
	}

	return append(result, StringValue{XMLName: xml.Name{Local: "string"}, Value: value})
}

type BooleanValue struct {
	XMLName xml.Name
}

func NewBooleanEntry(key string, value bool) []DictItem {
	return []DictItem{
		StringValue{XMLName: xml.Name{Local: "key"}, Value: key},
		BooleanValue{XMLName: xml.Name{Local: fmt.Sprintf("%t", value)}},
	}
}

type IntegerValue struct {
	XMLName xml.Name
	Value   int `xml:",chardata"`
}

func NewIntegerValue(key string, value int) []DictItem {
	return []DictItem{
		StringValue{XMLName: xml.Name{Local: "key"}, Value: key},
		IntegerValue{XMLName: xml.Name{Local: "integer"}, Value: value},
	}
}

func isSocketSSLorStartTLS(st providers.SocketType) bool {
	return st == providers.SSLSocketType || st == providers.StartTLSSocketType
}

type MobileConfigParams struct {
	Domain      string
	DisplayName string
	Username    string
	Provider    providers.Provider
}

func NewMobileConfig(p MobileConfigParams) (*MobileConfig, error) {
	if p.Provider.ID == "" {
		return nil, fmt.Errorf("invalid provider: ID is empty")
	}

	if p.Provider.ImapServer == nil {
		return nil, fmt.Errorf("no IMAP server configured for provider %s", p.Provider.ID)
	}

	if p.Provider.SmtpServer == nil {
		return nil, fmt.Errorf("no SMTP server configured for provider %s", p.Provider.ID)
	}
	payloadIdentifier := fmt.Sprintf("%s.mobileconfig", p.Domain)
	payloadDescription := fmt.Sprintf("Email account configuration for [%s]", p.DisplayName)

	arrayKeyValuePairs := [][]DictItem{
		NewStringEntry("EmailAccountDescription", p.DisplayName),
		NewStringEntry("EmailAccountName", p.DisplayName),
		NewStringEntry("EmailAccountType", "EmailTypeIMAP"),
		NewStringEntry("EmailAddress", p.Username),
		NewStringEntry("IncomingMailServerAuthentication", "EmailAuthPassword"),
		NewStringEntry("IncomingMailServerHostName", p.Provider.ImapServer.Hostname),
		NewIntegerValue("IncomingMailServerPortNumber", p.Provider.ImapServer.Port),
		NewBooleanEntry("IncomingMailServerUseSSL", isSocketSSLorStartTLS(p.Provider.ImapServer.SocketType)),
		NewStringEntry("IncomingMailServerUsername", p.Username),
		NewStringEntry("IncomingPassword", ""),
		NewStringEntry("OutgoingMailServerAuthentication", "EmailAuthPassword"),
		NewStringEntry("OutgoingMailServerHostName", p.Provider.SmtpServer.Hostname),
		NewIntegerValue("OutgoingMailServerPortNumber", p.Provider.SmtpServer.Port),
		NewBooleanEntry("OutgoingMailServerUseSSL", isSocketSSLorStartTLS(p.Provider.SmtpServer.SocketType)),
		NewStringEntry("OutgoingMailServerUsername", p.Username),
		NewBooleanEntry("OutgoingPasswordSameAsIncomingPassword", true),
		NewStringEntry("PayloadDescription", payloadDescription),
		NewStringEntry("PayloadDisplayName", p.DisplayName),
		NewStringEntry("PayloadIdentifier", payloadIdentifier),
		NewStringEntry("PayloadOrganization", p.Domain),
		NewStringEntry("PayloadType", payloadType),
		NewStringEntry("PayloadUUID", uuid.NewString()),
		NewIntegerValue("PayloadVersion", 1),
		// Options
		NewBooleanEntry("PreventAppSheet", false),
		NewBooleanEntry("PreventMove", false),
		NewBooleanEntry("SMIMEEnabled", false),
		NewBooleanEntry("allowMailDrop", true),
		NewBooleanEntry("SMIMEEnablePerMessageSwitch", false),
		NewBooleanEntry("SMIMESigningEnabled", false),
		NewBooleanEntry("disableMailRecentsSyncing", false),
	}

	var arrayContent []DictItem
	for _, keyValuePair := range arrayKeyValuePairs {
		arrayContent = append(arrayContent, keyValuePair...)
	}

	payloadContentKeyValuePairs := [][]DictItem{
		NewStringEntry("PayloadDescription", payloadDescription),
		NewStringEntry("PayloadDisplayName", p.DisplayName),
		NewStringEntry("PayloadIdentifier", payloadIdentifier),
		NewStringEntry("PayloadOrganization", p.Domain),
		NewBooleanEntry("PayloadRemovalDisallowed", false),
		NewStringEntry("PayloadType", "Configuration"),
		NewStringEntry("PayloadUUID", uuid.NewString()),
		NewIntegerValue("PayloadVersion", 1),
	}

	var payloadContent []DictItem
	for _, keyValuePair := range payloadContentKeyValuePairs {
		payloadContent = append(payloadContent, keyValuePair...)
	}

	return &MobileConfig{
		Version: plistVersion,
		Dict: Dict{
			Key:   "PayloadContent",
			Array: &Array{Items: []Dict{{Items: arrayContent}}},
			Items: payloadContent,
		},
	}, nil
}
