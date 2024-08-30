package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
		r.Post("/autodiscover/autodiscover.xml", handlers.AutodiscoverHandler)
		r.Post("/Autodiscover/Autodiscover.xml", handlers.AutodiscoverHandler)
	})

	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
