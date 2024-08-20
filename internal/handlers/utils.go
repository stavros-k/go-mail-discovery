package handlers

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/stavros-k/go-mail-discovery/internal/providers"
	"golang.org/x/net/publicsuffix"
)

var (
	ErrInvalidEmail = errors.New("invalid email address")
	ErrNoProvider   = errors.New("no provider found")
)

func handleError(w http.ResponseWriter, status int, err error) {
	log.Printf("Error: %v", err)
	http.Error(w, err.Error(), status)
}

func getEmailFromQuery(r *http.Request) (string, error) {
	emailAddress := r.URL.Query().Get("emailaddress")
	if emailAddress == "" {
		return "%EMAILADDRESS%", nil
	}
	if !strings.Contains(emailAddress, "@") {
		return "", ErrInvalidEmail
	}
	return emailAddress, nil
}

func getDomainFromRequest(r *http.Request) (string, error) {
	// TODO: Remove this line when ready for production
	r.Host = "escapegameover.be:8080"

	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		// TODO: Test this, and see if we can use error.Is
		// If SplitHostPort fails, it might be because there's no port
		// In this case, use the whole Host
		host = r.Host
	}

	domain, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return "", fmt.Errorf("error getting effective TLD: %w", err)
	}

	return domain, nil
}

func getProviderFromMX(domain string) (providers.Provider, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return providers.Provider{}, fmt.Errorf("error looking up MX records: %w", err)
	}

	for _, mxRecord := range mxRecords {
		switch {
		case strings.HasSuffix(mxRecord.Host, ".zoho.eu."):
			return providers.GetProvider("zoho,eu")
			// Add more cases here for other providers
		}
	}

	return providers.Provider{}, ErrNoProvider
}
