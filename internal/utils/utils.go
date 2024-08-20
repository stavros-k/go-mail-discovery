package utils

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/stavros-k/go-mail-discovery/internal/providers"
)

var (
	ErrInvalidEmail = errors.New("invalid email address")
	ErrNoProvider   = errors.New("no provider found")
)

func GetDomainFromEmailAddress(emailAddress string) (string, error) {
	if emailAddress == "" {
		return "", nil
	}
	if !strings.Contains(emailAddress, "@") {
		return "", ErrInvalidEmail
	}

	domain := strings.Split(emailAddress, "@")[1]

	if domain == "" {
		return "", ErrInvalidEmail
	}

	return domain, nil
}

func GetProviderFromMX(domain string) (providers.Provider, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return providers.Provider{}, fmt.Errorf("error looking up MX records: %w", err)
	}

	for _, mxRecord := range mxRecords {
		switch {
		case strings.HasSuffix(mxRecord.Host, ".zoho.eu."):
			return providers.GetProvider("zoho.eu")
			// Add more cases here for other providers
		}
	}

	return providers.Provider{}, ErrNoProvider
}
