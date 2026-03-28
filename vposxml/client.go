package vposxml

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/beevik/etree"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/internal/digest"
)

// Client performs VPOS XML 2.1 requests over HTTPS.
type Client struct {
	cfg    cardlink.Config
	client *http.Client
}

// NewClient returns a client using cfg for endpoint resolution and digest signing.
func NewClient(cfg cardlink.Config) *Client {
	hc := cfg.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{cfg: cfg, client: hc}
}

func (c *Client) postXML(ctx context.Context, body string) ([]byte, error) {
	u, err := c.cfg.VPOSXMLURL()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/xml")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return b, fmt.Errorf("vposxml: HTTP %s: %s", resp.Status, truncate(string(b), 500))
	}
	return b, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// verifyResponseDigest checks the response Message digest when present.
func verifyResponseDigest(xmlBytes []byte, secret string) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlBytes); err != nil {
		return err
	}
	root := doc.Root()
	if root == nil {
		return fmt.Errorf("vposxml: empty document")
	}
	var msgEl, digEl *etree.Element
	switch root.Tag {
	case "VPOS":
		for _, childEl := range root.ChildElements() {
			switch childEl.Tag {
			case "Message":
				msgEl = childEl
			case "Digest":
				digEl = childEl
			}
		}
	case "Message":
		msgEl = root
	}
	if msgEl == nil {
		return nil
	}
	// Error envelope uses Message version 1.0 — no standard digest verification
	if msgEl.SelectAttrValue("version", "") == "1.0" {
		return nil
	}
	if digEl == nil {
		return nil
	}
	got := strings.TrimSpace(digEl.Text())
	if got == "" {
		return nil
	}
	c14n, err := digest.CanonicalXML10Message(msgEl)
	if err != nil {
		return err
	}
	want := digest.VPOS21(c14n, secret)
	if got != want {
		return fmt.Errorf("vposxml: response digest mismatch")
	}
	return nil
}

// responseBytes returns raw XML after optional digest verification.
func (c *Client) responseBytes(ctx context.Context, body string, verify bool) ([]byte, error) {
	raw, err := c.postXML(ctx, body)
	if err != nil {
		return nil, err
	}
	if verify && c.cfg.SharedSecret != "" {
		if err := verifyResponseDigest(raw, c.cfg.SharedSecret); err != nil {
			return nil, err
		}
	}
	return raw, nil
}

// RawPOST sends pre-built VPOS XML (for advanced use).
func (c *Client) RawPOST(ctx context.Context, vposXML string, verifyResponse bool) ([]byte, error) {
	return c.responseBytes(ctx, vposXML, verifyResponse)
}
