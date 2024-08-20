package providers

import (
	"fmt"
	"sync"
)

var (
	providers      map[string]Provider
	providersMutex sync.RWMutex
)

func init() {
	providers = make(map[string]Provider)
	// Initialize static providers
	addProvider(zohoConfig())
	// Add more static providers here as needed
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
	for _, provider := range providers {
		fmt.Printf("%s\n", provider)
	}
	return nil
}
