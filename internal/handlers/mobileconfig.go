package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stavros-k/go-mail-discovery/internal/generators"
	"github.com/stavros-k/go-mail-discovery/internal/utils"
)

func MobileConfigHandler(w http.ResponseWriter, r *http.Request) {
	emailAddress, err := getEmailFromQuery(r, "email", false)
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

	config, err := generators.NewMobileConfig(generators.MobileConfigParams{
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

	w.Header().Set("Content-Type", "application/x-apple-aspen-config")
	// return a downloadable file with name "autoconfig.mobileconfig"
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.mobileconfig", emailAddress))
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Printf("error writing response: %v", err)
	}
}
