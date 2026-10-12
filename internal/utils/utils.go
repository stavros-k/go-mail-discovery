package utils

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/mail"
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
	// Only a bare address is accepted, not the "Name <address>" form
	address, err := mail.ParseAddress(emailAddress)
	if err != nil || address.Address != emailAddress {
		return "", ErrInvalidEmail
	}

	domain := NormalizeDomain(emailAddress[strings.LastIndex(emailAddress, "@")+1:])

	if domain == "" {
		return "", ErrInvalidEmail
	}

	return domain, nil
}

// NormalizeDomain lowercases the domain and strips the trailing dot,
// so that the same domain always maps to the same provider and cache entry.
func NormalizeDomain(domain string) string {
	return strings.ToLower(strings.TrimRight(domain, "."))
}

func init() {
	providerCacheMap = ProviderIDCache{
		cache: make(map[string]CachedProviderID),
	}
}

const (
	cacheTTL = time.Minute * 5
	// Domains without a known provider are cached for a shorter time,
	// so that repeated requests for them do not each cause a DNS lookup.
	negativeCacheTTL = time.Minute
	// Upper bound of cached domains, new entries are not cached once it is reached.
	maxCacheEntries = 10000
)

var providerCacheMap ProviderIDCache

// CachedProviderID holds the provider of a domain,
// an empty ProviderID means that the domain has no known provider.
type CachedProviderID struct {
	ProviderID string
	Expires    time.Time
}

type ProviderIDCache struct {
	cache map[string]CachedProviderID
	mutex sync.RWMutex
}

func (p *ProviderIDCache) Set(domain string, providerID string, ttl time.Duration) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	now := time.Now().UTC()
	for key, entry := range p.cache {
		if entry.Expires.Before(now) {
			delete(p.cache, key)
		}
	}
	if _, ok := p.cache[domain]; !ok && len(p.cache) >= maxCacheEntries {
		return
	}
	p.cache[domain] = CachedProviderID{
		ProviderID: providerID,
		Expires:    now.Add(ttl),
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
	domain = NormalizeDomain(domain)
	if cache {
		providerID, ok := providerCacheMap.Get(domain)
		if ok {
			if providerID == "" {
				return providers.Provider{}, ErrNoProvider
			}
			log.Printf("Using cached provider %s for domain %s", providerID, domain)
			return providers.GetProvider(providerID)
		}
	}

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		// Only cache definitive answers, not temporary failures like timeouts
		if dnsErr, ok := errors.AsType[*net.DNSError](err); ok && dnsErr.IsNotFound {
			providerCacheMap.Set(domain, "", negativeCacheTTL)
		}
		return providers.Provider{}, fmt.Errorf("error looking up MX records: %w", err)
	}

	for _, mxRecord := range mxRecords {
		host := NormalizeDomain(mxRecord.Host)
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
		providerCacheMap.Set(domain, provider.ID, cacheTTL)
		return provider, nil
	}

	providerCacheMap.Set(domain, "", negativeCacheTTL)
	return providers.Provider{}, ErrNoProvider
}
