package models

import (
	"ddns-updater/pkg/regex"
	"fmt"
	"time"
)

// SettingsType contains the elements to update the DNS record
type SettingsType struct {
	Domain      string
	Host        string
	Provider    ProviderType
	IPmethod    IPMethodType
	Delay       time.Duration
	NoDNSLookup bool
	// Provider dependent fields
	Password       string // Namecheap only
	Key            string // GoDaddy, Dreamhost and Cloudflare only
	Secret         string // GoDaddy only
	Token          string // DuckDNS only
	Email          string // Cloudflare only
	UserServiceKey string // Cloudflare only
	ZoneIdentifier string // Cloudflare only
	Identifier     string // Cloudflare only
	Proxied        bool   // Cloudflare only
}

func (settings *SettingsType) String() (s string) {
	return fmt.Sprintf("%s | %s | %s | %s", settings.Domain, settings.Host, settings.Provider, settings.IPmethod)
}

// BuildDomainName builds the domain name from the domain and the host of the settings
func (settings *SettingsType) BuildDomainName() string {
	if settings.Host == "@" {
		return settings.Domain
	} else if settings.Host == "*" {
		return settings.Domain // TODO random subdomain
	} else {
		return settings.Host + "." + settings.Domain
	}
}

func (settings *SettingsType) getHTMLDomain() string {
	return "<a href=\"http://" + settings.BuildDomainName() + "\">" + settings.Domain + "</a>"
}

func (settings *SettingsType) getHTMLProvider() string {
	switch settings.Provider {
	case PROVIDERNAMECHEAP:
		return "<a href=\"https://namecheap.com\">Namecheap</a>"
	case PROVIDERGODADDY:
		return "<a href=\"https://godaddy.com\">GoDaddy</a>"
	case PROVIDERDUCKDNS:
		return "<a href=\"https://duckdns.org\">DuckDNS</a>"
	case PROVIDERDREAMHOST:
		return "<a href=\"https://https://www.dreamhost.com/\">Dreamhost</a>"
	default:
		return settings.Provider.String()
	}
}

// TODO map to icons
func (settings *SettingsType) getHTMLIPMethod() string {
	switch settings.IPmethod {
	case IPMETHODPROVIDER:
		return settings.getHTMLProvider()
	case IPMETHODDUCKDUCKGO:
		return "<a href=\"https://duckduckgo.com/?q=ip\">DuckDuckGo</a>"
	case IPMETHODOPENDNS:
		return "<a href=\"https://diagnostic.opendns.com/myip\">OpenDNS</a>"
	default:
		return settings.IPmethod.String()
	}
}

// Verify verifies all the settings provided are valid
func (settings *SettingsType) Verify() error {
	if !regex.MatchDomain(settings.Domain) {
		return fmt.Errorf("the domain name %s is not valid for settings %s", settings.Domain, settings)
	}
	if len(settings.Host) == 0 {
		return fmt.Errorf("the host for entry %s must have at least one character", settings)
	}
	switch settings.Provider {
	case PROVIDERNAMECHEAP:
		if !regex.MatchNamecheapPassword(settings.Password) {
			return fmt.Errorf("the Namecheap password is not valid for settings %s", settings)
		}
	case PROVIDERGODADDY:
		if !regex.MatchGodaddyKey(settings.Key) {
			return fmt.Errorf("the GoDaddy key is not valid for settings %s", settings)
		}
		if !regex.MatchGodaddySecret(settings.Secret) {
			return fmt.Errorf("the GoDaddy secret is not valid for settings %s", settings)
		}
		if settings.IPmethod == IPMETHODPROVIDER {
			return fmt.Errorf("the provider %s does not support the IP update method %s", settings.Provider, settings.IPmethod)
		}
	case PROVIDERDUCKDNS:
		if !regex.MatchDuckDNSToken(settings.Token) {
			return fmt.Errorf("the DuckDNS token is not valid for settings %s", settings)
		}
		if settings.Host != "@" {
			return fmt.Errorf("the host %s can only be @ for settings %s", settings.Host, settings)
		}
	case PROVIDERDREAMHOST:
		if !regex.MatchDreamhostKey(settings.Key) {
			return fmt.Errorf("the Dreamhost key is not valid for settings %s", settings)
		}
		if settings.Host != "@" {
			return fmt.Errorf("the host %s can only be @ for settings %s", settings.Host, settings)
		}
		if settings.IPmethod == IPMETHODPROVIDER {
			return fmt.Errorf("the provider %s does not support the IP update method %s", settings.Provider, settings.IPmethod)
		}
	case PROVIDERCLOUDFLARE:
		if settings.UserServiceKey == "" { // email and key must be provided
			if !regex.MatchCloudflareKey(settings.Key) {
				return fmt.Errorf("the Cloudflare key is not valid for settings %s", settings)
			}
			if !regex.MatchEmail(settings.Email) {
				return fmt.Errorf("the Cloudflare email %s is not valid for settings %s", settings.Email, settings)
			}
		} else { // only user service key
			if !regex.MatchCloudflareUserServiceKey(settings.UserServiceKey) {
				return fmt.Errorf("the Cloudflare user service key is not valid for settings %s", settings)
			}
		}
		if len(settings.ZoneIdentifier) == 0 {
			return fmt.Errorf("Cloudflare zone identifier was not provided")
		}
		if len(settings.Identifier) == 0 {
			return fmt.Errorf("Cloudflare identifier was not provided")
		}
		if settings.IPmethod == IPMETHODPROVIDER {
			return fmt.Errorf("the provider %s does not support the IP update method %s", settings.Provider, settings.IPmethod)
		}
	default:
		return fmt.Errorf("provider %s is not supported", settings.Provider)
	}
	return nil
}
