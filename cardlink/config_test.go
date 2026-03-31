package cardlink

import (
	"strings"
	"testing"
)

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
	t.Run("should return sandbox host", func(t *testing.T) {
		cfg := Config{
			Environment: Sandbox,
			Partner:     Cardlink,
		}
		u, err := cfg.VPOSXMLURL()
		if err != nil {
			t.Fatal(err)
		}
		if u.Scheme != "https" || u.Host != "ecommerce-test.cardlink.gr" {
			t.Fatalf("got %s", u.String())
		}
	})

	t.Run("should return worldline production host", func(t *testing.T) {
		cfg := Config{
			Environment: Production,
			Partner:     Worldline,
		}
		u, err := cfg.VPOSXMLURL()
		if err != nil {
			t.Fatal(err)
		}
		if u.Scheme != "https" || u.Host != "eurocommerce.cardlink.gr" {
			t.Fatalf("got %s", u.String())
		}
	})
}
