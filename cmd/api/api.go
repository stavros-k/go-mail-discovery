package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stavros-k/go-mail-discovery/internal/handlers"
	"github.com/stavros-k/go-mail-discovery/internal/providers"
)

func main() {
	fmt.Println(strings.Join(providers.ListProvidersWithInfo(), "\n"))
	r := chi.NewRouter()

	r.Get("/health", handlers.HealthHandler)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/mail/config-v1.1.xml", handlers.AutoconfigHandler)
		r.Get("/email.mobileconfig", handlers.MobileConfigHandler)
		// Clients use different casings for this path
		r.Post("/{dir:(?i)autodiscover}/{file:(?i)autodiscover\\.xml}", handlers.AutodiscoverHandler)
	})

	fmt.Println("Starting server on port 8080")
	server := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}
	log.Fatal(server.ListenAndServe())
}
