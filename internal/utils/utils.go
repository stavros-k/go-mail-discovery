package utils

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/stavros-k/go-mail-discovery/internal/providers"
	"golang.org/x/net/publicsuffix"
)

var (
	ErrInvalidEmail = errors.New("invalid email address")
	ErrNoProvider   = errors.New("no provider found")
)

func GetDomainFromEmailAddress(emailAddress string) (string, error) {
	if emailAddress == "" {
		return "", nil
	}
	if !strings.Contains(emailAddress, "@") {
		return "", ErrInvalidEmail
	}

	domain := strings.Split(emailAddress, "@")[1]

	if domain == "" {
		return "", ErrInvalidEmail
	}

	return domain, nil
}

func init() {
	providerCacheMap = ProviderIDCache{
		cache: make(map[string]CachedProviderID),
	}
}

const cacheTTL = time.Minute * 5

var providerCacheMap ProviderIDCache

type CachedProviderID struct {
	ProviderID string
	Expires    time.Time
}

type ProviderIDCache struct {
	cache map[string]CachedProviderID
	mutex sync.RWMutex
}

func (p *ProviderIDCache) Set(domain string, providerID string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.cache[domain] = CachedProviderID{
		ProviderID: providerID,
		Expires:    time.Now().UTC().Add(cacheTTL),
	}
}

func (p *ProviderIDCache) Get(domain string) (string, bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	cacheEntry, ok := p.cache[domain]
	if !ok || cacheEntry.Expires.Before(time.Now().UTC()) {
		return "", false
	}
	return cacheEntry.ProviderID, ok
}

func GetProviderFromMX(domain string, cache bool) (providers.Provider, error) {
	if cache {
		providerID, ok := providerCacheMap.Get(domain)
		if ok {
			log.Printf("Using cached provider %s for domain %s", providerID, domain)
			return providers.GetProvider(providerID)
		}
	}

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return providers.Provider{}, fmt.Errorf("error looking up MX records: %w", err)
	}

	for _, mxRecord := range mxRecords {
		host := strings.TrimRight(mxRecord.Host, ".")
		if strings.Contains(host, ":") {
			host, _, err = net.SplitHostPort(host)
			if err != nil {
				log.Printf("error parsing URL: %v", err)
				continue
			}
		}
		parsedDomain, err := publicsuffix.EffectiveTLDPlusOne(host)
		if err != nil {
			log.Printf("error getting effective TLD: %v", err)
			continue
		}
		provider, err := providers.GetProvider(parsedDomain)
		if err != nil {
			continue
		}

		log.Printf("Updating cache provider %s for domain %s", provider.ID, domain)
		providerCacheMap.Set(domain, provider.ID)
		return provider, nil
	}

	return providers.Provider{}, ErrNoProvider
}
