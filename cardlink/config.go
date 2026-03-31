package cardlink

import (
	"fmt"
	"net/url"
)

const (
	vposXMLPath = "/vpos/xmlpayvpos"
)

// Config holds merchant credentials and endpoint resolution for VPOS XML.
type Config struct {
	MID          string
	SharedSecret string
	Environment  Environment
	Partner      BusinessPartner
}

func sandboxHosts() map[BusinessPartner]string {
	return map[BusinessPartner]string{
		Cardlink:  "https://ecommerce-test.cardlink.gr",
		Nexi:      "https://alphaecommerce-test.cardlink.gr",
		Worldline: "https://eurocommerce-test.cardlink.gr",
	}
}

func productionHosts() map[BusinessPartner]string {
	return map[BusinessPartner]string{
		Cardlink:  "https://ecommerce.cardlink.gr",
		Nexi:      "https://alphaecommerce.cardlink.gr",
		Worldline: "https://eurocommerce.cardlink.gr",
	}
}

// resolveBase returns the gateway origin for VPOS XML requests, without trailing slash.
func (c *Config) resolveBase() string {
	if c.Environment == Sandbox {
		return sandboxHosts()[c.Partner]
	}
	return productionHosts()[c.Partner]
}

// VPOSXMLURL returns the full URL for VPOS XML 2.1 requests.
func (c *Config) VPOSXMLURL() (*url.URL, error) {
	base := c.resolveBase()
	if base == "" {
		return nil, fmt.Errorf("cardlink: unknown BusinessPartner %v", c.Partner)
	}

	return url.Parse(base + vposXMLPath)
}
