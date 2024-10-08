package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stavros-k/go-mail-discovery/internal/generators"
	"github.com/stavros-k/go-mail-discovery/internal/utils"
)

func AutoconfigHandler(w http.ResponseWriter, r *http.Request) {
	emailAddress, err := getEmailFromQuery(r, "emailaddress", true)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("error getting email from query: %w", err))
		return
	}

	domain, err := getDomain(emailAddress, r)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("error getting domain: %w", err))
		return
	}
	cache := r.URL.Query().Get("cache") != "false"
	provider, err := utils.GetProviderFromMX(domain, cache)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("error getting provider from MX: %w", err))
		return
	}

	config, err := generators.NewConfigV1_1(generators.ConfigV1_1Params{
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
