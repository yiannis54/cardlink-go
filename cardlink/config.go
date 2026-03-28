package cardlink

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	redirectPath = "/vpos/shophandlermpi"
	vposXMLPath  = "/vpos/xmlpayvpos"
)

// Config holds merchant credentials and endpoint resolution for redirect and VPOS XML.
type Config struct {
	MID          string
	SharedSecret string
	Environment  Environment
	Partner      BusinessPartner

	// RedirectBaseURL overrides the default gateway origin (no trailing slash), e.g. "https://ecommerce-test.cardlink.gr"
	RedirectBaseURL string
	// VPOSXMLBaseURL overrides the default origin for POST /vpos/xmlpayvpos
	VPOSXMLBaseURL string

	// HTTPClient is optional; vposxml.Client uses http.DefaultClient if nil
	HTTPClient *http.Client
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

// resolveBase returns the gateway origin for redirect and XML, without trailing slash.
func (c *Config) resolveBase() string {
	if c.Environment == Sandbox {
		return sandboxHosts()[c.Partner]
	}
	return productionHosts()[c.Partner]
}

// RedirectURL returns the full POST URL for redirect payment initiation.
func (c *Config) RedirectURL() (*url.URL, error) {
	base := strings.TrimSuffix(c.RedirectBaseURL, "/")
	if base == "" {
		base = c.resolveBase()
	}
	if base == "" {
		return nil, fmt.Errorf("cardlink: unknown BusinessPartner %s", c.Partner)
	}
	u, err := url.Parse(base + redirectPath)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// VPOSXMLURL returns the full URL for VPOS XML 2.1 requests.
func (c *Config) VPOSXMLURL() (*url.URL, error) {
	base := strings.TrimSuffix(c.VPOSXMLBaseURL, "/")
	if base == "" {
		base = c.resolveBase()
	}
	if base == "" {
		return nil, fmt.Errorf("cardlink: unknown BusinessPartner %v", c.Partner)
	}
	u, err := url.Parse(base + vposXMLPath)
	if err != nil {
		return nil, err
	}
	return u, nil
}
