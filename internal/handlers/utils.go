package handlers

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/stavros-k/go-mail-discovery/internal/utils"
	"golang.org/x/net/publicsuffix"
)

var (
	ErrInvalidEmail = errors.New("invalid email address")
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
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		if strings.Contains(host, ":") {
			return "", fmt.Errorf("error splitting host and port: %w", err)
		}
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

func getDomain(emailAddress string, r *http.Request) (string, error) {
	if emailAddress != "" && emailAddress != "%EMAILADDRESS%" {
		return utils.GetDomainFromEmailAddress(emailAddress)
	}
	return getDomainFromRequest(r)
}
