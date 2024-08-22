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
	r.Use(middleware.Logger)
	r.Get("/mail/config-v1.1.xml", handlers.AutoconfigHandler)
	r.Get("/test", handlers.MobileConfigHandler)
	// r.Post("/autodiscover/autodiscover.xml", autodiscoverHandler)
	// r.Post("/Autodiscover/Autodiscover.xml", autodiscoverHandler)

	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
