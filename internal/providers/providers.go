package providers

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/goccy/go-yaml"
)

var (
	providers      map[string]Provider
	providersMutex sync.RWMutex
)

func init() {
	providers = make(map[string]Provider)
	configPath := os.Getenv("MAIL_PROVIDER_CONFIG_PATH")
	if configPath == "" {
		log.Printf("Env MAIL_PROVIDER_CONFIG_PATH is not set, using default providers")
		if err := addProvider(zohoConfig()); err != nil {
			log.Fatalf("error adding default provider: %v", err)
		}
	} else {
		if err := LoadFromYaml(configPath); err != nil {
			log.Fatalf("error loading providers from config: %v", err)
		}
	}

	for _, provider := range providers {
		if err := provider.Validate(); err != nil {
			log.Fatalf("provider %s validation failed: %v", provider.ID, err)
		}
	}

}

// normalizeID lowercases the ID, as it is matched against
// the (lowercased) domain of the MX records.
func normalizeID(id string) string {
	return strings.ToLower(id)
}

func addProvider(provider Provider) error {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	provider.ID = normalizeID(provider.ID)
	if _, ok := providers[provider.ID]; ok {
		return fmt.Errorf("provider already exists: %s", provider.ID)
	}
	providers[provider.ID] = provider
	return nil
}

func GetProvider(id string) (Provider, error) {
	providersMutex.RLock()
	defer providersMutex.RUnlock()
	provider, ok := providers[normalizeID(id)]
	if !ok {
		return Provider{}, fmt.Errorf("provider not found: %s", id)
	}
	return provider, nil
}

func ListProviders() []string {
	providersMutex.RLock()
	defer providersMutex.RUnlock()
	ids := make([]string, 0, len(providers))
	for id := range providers {
		ids = append(ids, id)
	}
	return ids
}

func ListProvidersWithInfo() []string {
	providersMutex.RLock()
	defer providersMutex.RUnlock()
	providerInfos := make([]string, 0, len(providers))
	for _, provider := range providers {
		providerInfos = append(providerInfos, provider.String())
	}
	return providerInfos
}

func GetProviderInfo(id string) (string, error) {
	providersMutex.RLock()
	defer providersMutex.RUnlock()
	provider, ok := providers[normalizeID(id)]
	if !ok {
		return "", fmt.Errorf("provider not found: %s", id)
	}
	return provider.String(), nil
}

func LoadFromYaml(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var providers Providers
	err = yaml.Unmarshal(data, &providers)
	if err != nil {
		return err
	}
	for _, provider := range providers.Providers {
		if err := addProvider(provider); err != nil {
			return err
		}
	}
	return nil
}
