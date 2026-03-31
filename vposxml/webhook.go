package vposxml

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/beevik/etree"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/internal/digest"
)

// WebhookResponse holds fields parsed from a VPOS XML notification
// (e.g. SaleResponse, AuthorisationResponse, or Advice after a payment link payment).
type WebhookResponse struct {
	OrderID      string
	OrderAmount  string
	Currency     string
	PaymentTotal string
	Status       cardlink.Status
	TxID         string
	PaymentRef   string
	RiskScore    string
	ErrorCode    string
	Description  string
	RawXML       []byte
}

// VerifyWebhook verifies the VPOS XML 2.1 digest on a webhook/notification
// XML body and parses the response fields.
//
// Use this when you receive notifications/callbacks from VPOS XML payment
// links. The XML body is expected to contain:
//
//	<VPOS><Message ...>...</Message><Digest>...</Digest></VPOS>
func (c *Client) VerifyWebhook(xmlBody []byte) (*WebhookResponse, error) {
	if c.cfg.SharedSecret == "" {
		return nil, ErrMissingSecret
	}
	if err := verifyWebhookDigest(xmlBody, c.cfg.SharedSecret); err != nil {
		return nil, err
	}
	return parseWebhookResponse(xmlBody)
}

// VerifyWebhookRequest reads the body of r and calls VerifyWebhook.
func (c *Client) VerifyWebhookRequest(r *http.Request) (*WebhookResponse, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("vposxml: reading request body: %w", err)
	}
	return c.VerifyWebhook(body)
}

// verifyWebhookDigest is stricter than verifyResponseDigest:
// it requires the Digest element to be present and non-empty (except for v1.0 error envelopes).
func verifyWebhookDigest(xmlBytes []byte, secret string) error {
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
		return fmt.Errorf("vposxml: no Message element in webhook payload")
	}
	if msgEl.SelectAttrValue("version", "") == "1.0" {
		return nil
	}
	if digEl == nil {
		return ErrMissingDigest
	}
	got := strings.TrimSpace(digEl.Text())
	if got == "" {
		return ErrMissingDigest
	}

	c14n, err := digest.CanonicalXML10Message(msgEl)
	if err != nil {
		return err
	}
	want := digest.VPOS21(c14n, secret)
	if got != want {
		return ErrDigestMismatch
	}
	return nil
}

func parseWebhookResponse(raw []byte) (*WebhookResponse, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(raw); err != nil {
		return nil, err
	}
	root := doc.Root()

	var msg *etree.Element
	if root != nil && root.Tag == "VPOS" {
		msg = root.SelectElement("Message")
	} else if root != nil && root.Tag == "Message" {
		msg = root
	}
	if msg == nil {
		return &WebhookResponse{RawXML: raw}, nil
	}
	if re := responseErrorFromMessage(msg); re != nil {
		return nil, re
	}

	var resp *etree.Element
	for _, child := range msg.ChildElements() {
		tag := child.Tag
		if strings.HasSuffix(tag, "Response") || strings.HasSuffix(tag, "Notification") || tag == "Advice" {
			resp = child
			break
		}
	}
	if resp == nil {
		return &WebhookResponse{RawXML: raw}, nil
	}

	statusStr := textOrEmpty(resp, "./Status")
	if statusStr == "" {
		statusStr = firstText(resp, "TxStatus", "OrderTxStatus")
	}

	return &WebhookResponse{
		OrderID:      textOrEmpty(resp, "./OrderId"),
		OrderAmount:  textOrEmpty(resp, "./OrderAmount"),
		Currency:     textOrEmpty(resp, "./Currency"),
		PaymentTotal: textOrEmpty(resp, "./PaymentTotal"),
		Status:       cardlink.ParseStatus(statusStr),
		TxID:         firstText(resp, "TxId"),
		PaymentRef:   textOrEmpty(resp, "./PaymentRef"),
		RiskScore:    textOrEmpty(resp, "./RiskScore"),
		ErrorCode:    textOrEmpty(resp, "./ErrorCode"),
		Description:  textOrEmpty(resp, "./Description"),
		RawXML:       raw,
	}, nil
}
