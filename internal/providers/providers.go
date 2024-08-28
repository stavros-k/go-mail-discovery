package providers

import (
	"fmt"
	"io"
	"log"
	"os"
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
		addProvider(zohoConfig())
	} else {
		if err := LoadFromYaml(configPath); err != nil {
			log.Printf("error loading providers from config: %v", err)
		}
	}

	for _, provider := range providers {
		if err := provider.Validate(); err != nil {
			log.Fatalf("provider %s validation failed: %v", provider.ID, err)
		}
	}

}

func addProvider(provider Provider) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	if _, ok := providers[provider.ID]; ok {
		panic(fmt.Sprintf("provider already exists: %s", provider.ID))
	}
	providers[provider.ID] = provider
}

func GetProvider(id string) (Provider, error) {
	providersMutex.RLock()
	defer providersMutex.RUnlock()
	provider, ok := providers[id]
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
	provider, ok := providers[id]
	if !ok {
		return "", fmt.Errorf("provider not found: %s", id)
	}
	return provider.String(), nil
}

func LoadFromYaml(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var providers Providers
	err = yaml.Unmarshal(data, &providers)
	if err != nil {
		return err
	}
	for _, provider := range providers.Providers {
		addProvider(provider)
	}
	return nil
}
