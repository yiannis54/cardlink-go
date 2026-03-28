package cardlink

import (
	"strings"
	"testing"
)

func TestConfig_RedirectURL(t *testing.T) {
	cfg := Config{Environment: Sandbox, Partner: Worldline}
	u, err := cfg.RedirectURL()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(u.Path, "/vpos/shophandlermpi") {
		t.Fatalf("path: %s", u.Path)
	}
	if !strings.Contains(u.Host, "eurocommerce-test") {
		t.Fatalf("host: %s", u.Host)
	}
}

func TestConfig_VPOSXMLURL(t *testing.T) {
	cfg := Config{Environment: Production, Partner: Cardlink}
	u, err := cfg.VPOSXMLURL()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(u.Path, "/vpos/xmlpayvpos") {
		t.Fatalf("path: %s", u.Path)
	}
}

func TestConfig_OverrideBase(t *testing.T) {
	cfg := Config{
		Environment:     Sandbox,
		Partner:         Cardlink,
		RedirectBaseURL: "https://custom.example",
	}
	u, err := cfg.RedirectURL()
	if err != nil {
		t.Fatal(err)
	}
	if u.Scheme != "https" || u.Host != "custom.example" {
		t.Fatalf("got %s", u.String())
	}
}
