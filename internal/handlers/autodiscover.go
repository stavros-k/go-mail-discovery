package handlers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stavros-k/go-mail-discovery/internal/generators"
	"github.com/stavros-k/go-mail-discovery/internal/utils"
)

const maxAutodiscoverBodySize = 64 * 1024

type AutoDiscoverPayload struct {
	XMLName xml.Name `xml:"Autodiscover"`
	Request Request  `xml:"Request"`
}
type Request struct {
	XMLName      xml.Name `xml:"Request"`
	EmailAddress string   `xml:"EMailAddress"`
}

func AutodiscoverHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxAutodiscoverBodySize))
	if err != nil {
		status := http.StatusBadRequest
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			status = http.StatusRequestEntityTooLarge
		}
		handleError(w, status, fmt.Errorf("error reading request body: %w", err))
		return
	}
	var payload AutoDiscoverPayload
	err = xml.Unmarshal(body, &payload)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("error unmarshalling request body: %w", err))
		return
	}
	emailAddress := payload.Request.EmailAddress
	if emailAddress == "" || !strings.Contains(emailAddress, "@") {
		handleError(w, http.StatusBadRequest, fmt.Errorf("invalid email address: %s", emailAddress))
		return
	}

	domain, err := getDomain(emailAddress, r)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("error getting domain: %w", err))
		return
	}
	provider, err := utils.GetProviderFromMX(domain, true)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("error getting provider from MX: %w", err))
		return
	}

	config, err := generators.NewAutoDiscoverConfig(generators.AutoDiscoverConfigParams{
		Domain:      domain,
		DisplayName: emailAddress,
		Username:    emailAddress,
		Provider:    provider,
	})
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Errorf("error generating config: %w", err))
		return
	}
	data, err := config.Bytes()
	if err != nil {
		handleError(w, http.StatusInternalServerError, fmt.Errorf("error generating config: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Printf("error writing response: %v", err)
	}
}
