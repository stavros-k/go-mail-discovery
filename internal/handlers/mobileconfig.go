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

	provider, err := utils.GetProviderFromMX(domain)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("error getting provider from MX: %w", err))
		return
	}

	config, err := generators.NewMobileConfig(generators.MobileConfigParams{
		Domain:      "stavrosk.me",
		DisplayName: "stavros@stavrosk.me",
		Username:    "stavros@stavrosk.me",
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
