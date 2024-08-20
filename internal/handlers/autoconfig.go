package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stavros-k/go-mail-discovery/internal/generators"
)

func AutoconfigHandler(w http.ResponseWriter, r *http.Request) {
	emailAddress, err := getEmailFromQuery(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("error getting email from query: %w", err))
		return
	}

	domain, err := getDomainFromRequest(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("error getting domain from request: %w", err))
		return
	}

	provider, err := getProviderFromMX(domain)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("error getting provider from MX: %w", err))
		return
	}

	config := generators.NewConfig_v1_1_xml(generators.Config_v1_1_xml_params{
		Domain:      domain,
		DisplayName: emailAddress,
		Username:    emailAddress,
		Provider:    provider,
	})

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
