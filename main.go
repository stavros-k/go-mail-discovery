package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/stavros-k/go-mail-discovery/internal/generators"
	"github.com/stavros-k/go-mail-discovery/internal/providers"
	"golang.org/x/net/publicsuffix"
)

func main() {
	r := chi.NewRouter()
	r.Get("/mail/config-v1.1.xml", autoconfigHandler)
	// r.Post("/autodiscover/autodiscover.xml", autodiscoverHandler)
	// r.Post("/Autodiscover/Autodiscover.xml", autodiscoverHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func getEmailFromQuery(r *http.Request) (string, error) {
	emailaddress := r.URL.Query().Get("emailaddress")
	if emailaddress != "" {
		if !strings.Contains(emailaddress, "@") {
			return "", errors.New("invalid email address")
		}
		return emailaddress, nil
	}
	return "%EMAILADDRESS%", nil
}

func getDomainFromRequest(r *http.Request) (string, error) {
	r.Host = "escapegameover.be:8080"
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		return "", err
	}
	// Use the publicsuffix package to get the registered domain
	domain, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return "", err
	}

	return domain, nil
}

func getProviderFromMX(domain string) (providers.Server, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return providers.Server{}, err
	}
	for _, mxRecord := range mxRecords {
		switch {
		case strings.HasSuffix(mxRecord.Host, ".zoho.eu."):
			return providers.Zoho(), nil
		}
	}
	return providers.Server{}, errors.New("no provider found")
}

func autoconfigHandler(w http.ResponseWriter, r *http.Request) {
	emailaddress, err := getEmailFromQuery(r)
	if err != nil {
		log.Printf("error getting email from query: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	domain, err := getDomainFromRequest(r)
	if err != nil {
		log.Printf("error getting domain from request: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	provider, err := getProviderFromMX(domain)
	if err != nil {
		log.Printf("error getting provider from MX: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	config := generators.NewConfig_v1_1_xml(generators.Config_v1_1_xml_params{
		Domain:      domain,
		DisplayName: emailaddress,
		Username:    emailaddress,
		Provider:    provider,
	})

	data, err := config.Bytes()
	if err != nil {
		log.Printf("error generating config: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
