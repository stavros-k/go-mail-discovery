package generators

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"sync"
	"text/template"

	"github.com/google/uuid"
	"github.com/stavros-k/go-mail-discovery/internal/providers"
)

const (
	mobileConfigTemplateName = "mobileConfig"
	appleMailAppID           = "com.apple.mail.managed"
)

var (
	templates map[string]*template.Template
	mutex     sync.RWMutex
)

//go:embed mobileconfig.go.tmpl
var templateFiles embed.FS

func init() {
	// Parse templates on startup
	parsedTemplates, err := template.New(mobileConfigTemplateName).
		Funcs(template.FuncMap{
			"isSocketSSL": isSocketSSL,
		}).
		ParseFS(templateFiles, "mobileconfig.go.tmpl")
	if err != nil {
		log.Fatalf("error parsing templates: %v", err)
	}

	// Store the parsed templates in the map
	mutex.Lock()
	defer mutex.Unlock()
	templates = make(map[string]*template.Template)
	templates[mobileConfigTemplateName] = parsedTemplates
}

type MobileConfig struct {
	p                  MobileConfigParams
	uuid               string
	payloadDescription string
	payloadIdentifier  string
}

func isSocketSSL(st providers.SocketType) bool {
	return st == providers.SSLSocketType
}

func (m *MobileConfig) Bytes() ([]byte, error) {
	mutex.RLock()
	t, ok := templates[mobileConfigTemplateName]
	mutex.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no template found for provider %s", m.p.Provider.ID)
	}

	var b bytes.Buffer
	if err := t.ExecuteTemplate(&b, mobileConfigTemplateName, map[string]any{
		"AppleMailAppID":     appleMailAppID,
		"UUID":               m.uuid,
		"PayloadDescription": m.payloadDescription,
		"PayloadIdentifier":  m.payloadIdentifier,
		"ProviderID":         m.p.Provider.ID,
		"DisplayName":        m.p.DisplayName,
		"Username":           m.p.Username,
		"IMAPServer":         m.p.Provider.ImapServer,
		"SMTPServer":         m.p.Provider.SmtpServer,
	}); err != nil {
		return nil, fmt.Errorf("error executing template: %w", err)
	}
	return b.Bytes(), nil
}

func (m *MobileConfig) String() string {
	b, err := m.Bytes()
	if err != nil {
		return fmt.Sprintf("Error generating config: %v", err)
	}
	return string(b)
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

	uuid := uuid.NewString()
	payloadIdentifier := fmt.Sprintf("%s.autoconfig.%s", p.Provider.ID, uuid)
	payloadDescription := fmt.Sprintf("Email account configuration for [%s]", p.DisplayName)

	return &MobileConfig{
		p:                  p,
		uuid:               uuid,
		payloadDescription: payloadDescription,
		payloadIdentifier:  payloadIdentifier,
	}, nil
}
