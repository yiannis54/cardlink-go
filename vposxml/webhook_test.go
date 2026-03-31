package vposxml

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/beevik/etree"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/internal/digest"
)

const testSecret = "SecRetDigest1"

// buildTestWebhookXML constructs a VPOS XML envelope with a correctly computed digest.
func buildTestWebhookXML(t *testing.T, secret string) []byte {
	t.Helper()
	msgXML := `<Message xmlns="http://www.modirum.com/schemas/vposxmlapi41" xmlns:ns2="http://www.w3.org/2000/09/xmldsig#" messageId="M100" timeStamp="2025-03-28T12:00:00.000+02:00" version="2.1"><SaleResponse><OrderId>ORD123</OrderId><OrderAmount>10.00</OrderAmount><Currency>EUR</Currency><PaymentTotal>10.00</PaymentTotal><Status>CAPTURED</Status><TxId>999</TxId><PaymentRef>REF001</PaymentRef><RiskScore>5</RiskScore><Description>OK</Description></SaleResponse></Message>`

	msgEl, err := digest.ParseMessageElement(msgXML)
	if err != nil {
		t.Fatal(err)
	}
	c14n, err := digest.CanonicalXML10Message(msgEl)
	if err != nil {
		t.Fatal(err)
	}
	dig := digest.VPOS21(c14n, secret)

	doc := etree.NewDocument()
	doc.WriteSettings = etree.WriteSettings{CanonicalText: true}
	root := etree.NewElement("VPOS")
	root.CreateAttr("xmlns", vposNS)
	root.CreateAttr("xmlns:ns2", dsigNS)
	root.AddChild(msgEl)
	digEl := etree.NewElement("Digest")
	digEl.SetText(dig)
	root.AddChild(digEl)
	doc.SetRoot(root)

	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestVerifyWebhook_ValidDigest(t *testing.T) {
	cfg := cardlink.Config{MID: "0000001", SharedSecret: testSecret}
	c := NewClient(cfg)

	xml := buildTestWebhookXML(t, testSecret)
	resp, err := c.VerifyWebhook(xml)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OrderID != "ORD123" {
		t.Errorf("OrderID = %q, want ORD123", resp.OrderID)
	}
	if resp.Status != cardlink.StatusCaptured {
		t.Errorf("Status = %q, want CAPTURED", resp.Status)
	}
	if resp.TxID != "999" {
		t.Errorf("TxID = %q, want 999", resp.TxID)
	}
	if resp.PaymentRef != "REF001" {
		t.Errorf("PaymentRef = %q", resp.PaymentRef)
	}
	if resp.RiskScore != "5" {
		t.Errorf("RiskScore = %q", resp.RiskScore)
	}
	if resp.OrderAmount != "10.00" {
		t.Errorf("OrderAmount = %q", resp.OrderAmount)
	}
	if resp.Currency != "EUR" {
		t.Errorf("Currency = %q", resp.Currency)
	}
	if resp.PaymentTotal != "10.00" {
		t.Errorf("PaymentTotal = %q", resp.PaymentTotal)
	}
}

func TestVerifyWebhook_WrongSecret(t *testing.T) {
	cfg := cardlink.Config{MID: "0000001", SharedSecret: "WrongSecret"}
	c := NewClient(cfg)

	xml := buildTestWebhookXML(t, testSecret)
	_, err := c.VerifyWebhook(xml)
	if !errors.Is(err, ErrDigestMismatch) {
		t.Fatalf("expected ErrDigestMismatch, got %v", err)
	}
}

func TestVerifyWebhook_MissingDigestElement(t *testing.T) {
	cfg := cardlink.Config{MID: "0000001", SharedSecret: testSecret}
	c := NewClient(cfg)

	raw := []byte(`<VPOS xmlns="http://www.modirum.com/schemas/vposxmlapi41"><Message messageId="M1" timeStamp="2025-01-01T00:00:00+02:00" version="2.1"><SaleResponse><OrderId>O1</OrderId><Status>CAPTURED</Status></SaleResponse></Message></VPOS>`)
	_, err := c.VerifyWebhook(raw)
	if !errors.Is(err, ErrMissingDigest) {
		t.Fatalf("expected ErrMissingDigest, got %v", err)
	}
}

func TestVerifyWebhook_EmptyDigest(t *testing.T) {
	cfg := cardlink.Config{MID: "0000001", SharedSecret: testSecret}
	c := NewClient(cfg)

	raw := []byte(`<VPOS xmlns="http://www.modirum.com/schemas/vposxmlapi41"><Message messageId="M1" timeStamp="2025-01-01T00:00:00+02:00" version="2.1"><SaleResponse><OrderId>O1</OrderId><Status>CAPTURED</Status></SaleResponse></Message><Digest></Digest></VPOS>`)
	_, err := c.VerifyWebhook(raw)
	if !errors.Is(err, ErrMissingDigest) {
		t.Fatalf("expected ErrMissingDigest, got %v", err)
	}
}

func TestVerifyWebhook_MissingSecret(t *testing.T) {
	cfg := cardlink.Config{MID: "0000001", SharedSecret: ""}
	c := NewClient(cfg)

	_, err := c.VerifyWebhook([]byte(`<VPOS/>`))
	if !errors.Is(err, ErrMissingSecret) {
		t.Fatalf("expected ErrMissingSecret, got %v", err)
	}
}

func TestVerifyWebhook_ErrorEnvelope(t *testing.T) {
	cfg := cardlink.Config{MID: "0000001", SharedSecret: testSecret}
	c := NewClient(cfg)

	raw := []byte(`<VPOS xmlns="http://www.modirum.com/schemas/vposxmlapi41"><Message version="1.0" messageId="M1" timeStamp="2025-01-01T00:00:00+02:00"><ErrorCode>XE</ErrorCode><ErrorMessage>bad xml</ErrorMessage><Description>detail</Description></Message></VPOS>`)
	_, err := c.VerifyWebhook(raw)
	if err == nil {
		t.Fatal("expected error for error envelope")
	}
	var re *cardlink.ResponseError
	if !errors.As(err, &re) {
		t.Fatalf("expected ResponseError, got %v", err)
	}
	if re.ErrorCode != "XE" {
		t.Errorf("ErrorCode = %q", re.ErrorCode)
	}
}

func TestVerifyWebhookRequest(t *testing.T) {
	cfg := cardlink.Config{MID: "0000001", SharedSecret: testSecret}
	c := NewClient(cfg)

	xml := buildTestWebhookXML(t, testSecret)
	req, _ := http.NewRequest(http.MethodPost, "/webhook", io.NopCloser(bytes.NewReader(xml)))
	resp, err := c.VerifyWebhookRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OrderID != "ORD123" {
		t.Errorf("OrderID = %q", resp.OrderID)
	}
}

func TestVerifyWebhook_AuthorisationResponse(t *testing.T) {
	msgXML := `<Message xmlns="http://www.modirum.com/schemas/vposxmlapi41" xmlns:ns2="http://www.w3.org/2000/09/xmldsig#" messageId="M200" timeStamp="2025-03-28T12:00:00.000+02:00" version="2.1"><AuthorisationResponse><OrderId>PRE1</OrderId><OrderAmount>50.00</OrderAmount><Currency>EUR</Currency><PaymentTotal>50.00</PaymentTotal><Status>AUTHORIZED</Status><TxId>888</TxId><PaymentRef>REF002</PaymentRef></AuthorisationResponse></Message>`

	msgEl, err := digest.ParseMessageElement(msgXML)
	if err != nil {
		t.Fatal(err)
	}
	c14n, err := digest.CanonicalXML10Message(msgEl)
	if err != nil {
		t.Fatal(err)
	}
	dig := digest.VPOS21(c14n, testSecret)

	doc := etree.NewDocument()
	doc.WriteSettings = etree.WriteSettings{CanonicalText: true}
	root := etree.NewElement("VPOS")
	root.CreateAttr("xmlns", vposNS)
	root.CreateAttr("xmlns:ns2", dsigNS)
	root.AddChild(msgEl)
	digEl := etree.NewElement("Digest")
	digEl.SetText(dig)
	root.AddChild(digEl)
	doc.SetRoot(root)

	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		t.Fatal(err)
	}

	cfg := cardlink.Config{MID: "0000001", SharedSecret: testSecret}
	c := NewClient(cfg)

	resp, err := c.VerifyWebhook(buf.Bytes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != cardlink.StatusAuthorized {
		t.Errorf("Status = %q, want AUTHORIZED", resp.Status)
	}
	if resp.TxID != "888" {
		t.Errorf("TxID = %q", resp.TxID)
	}
}

func TestVerifyWebhook_Advice(t *testing.T) {
	msgXML := `<Message xmlns="http://www.modirum.com/schemas/vposxmlapi41" xmlns:ns2="http://www.w3.org/2000/09/xmldsig#" messageId="ADV12345678912345" timeStamp="2026-03-29T01:44:29.071+02:00" version="2.1"><Advice type="Sale"><Authentication><Mid>9000011260</Mid></Authentication><OrderId>123465</OrderId><OrderAmount>446.4</OrderAmount><Currency>EUR</Currency><OrderTxId>12345678912345</OrderTxId><OrderTxStatus>ERROR</OrderTxStatus><PaymentTotal>446.4</PaymentTotal><TxId>12345678912345</TxId><TxStatus>ERROR</TxStatus><TxTotal>446.4</TxTotal><TxCurrency>EUR</TxCurrency></Advice></Message>`

	msgEl, err := digest.ParseMessageElement(msgXML)
	if err != nil {
		t.Fatal(err)
	}
	c14n, err := digest.CanonicalXML10Message(msgEl)
	if err != nil {
		t.Fatal(err)
	}
	dig := digest.VPOS21(c14n, testSecret)

	doc := etree.NewDocument()
	doc.WriteSettings = etree.WriteSettings{CanonicalText: true}
	root := etree.NewElement("VPOS")
	root.CreateAttr("xmlns", vposNS)
	root.CreateAttr("xmlns:ns2", dsigNS)
	root.AddChild(msgEl)
	digEl := etree.NewElement("Digest")
	digEl.SetText(dig)
	root.AddChild(digEl)
	doc.SetRoot(root)

	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		t.Fatal(err)
	}

	cfg := cardlink.Config{MID: "0000001", SharedSecret: testSecret}
	c := NewClient(cfg)

	resp, err := c.VerifyWebhook(buf.Bytes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OrderID != "123465" {
		t.Errorf("OrderID = %q, want 123465", resp.OrderID)
	}
	if resp.Status != cardlink.StatusError {
		t.Errorf("Status = %q, want ERROR", resp.Status)
	}
	if resp.TxID != "12345678912345" {
		t.Errorf("TxID = %q", resp.TxID)
	}
	if resp.OrderAmount != "446.4" || resp.Currency != "EUR" || resp.PaymentTotal != "446.4" {
		t.Errorf("amounts: OrderAmount=%q Currency=%q PaymentTotal=%q", resp.OrderAmount, resp.Currency, resp.PaymentTotal)
	}
}

func TestParseWebhookResponse(t *testing.T) {
	raw := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><VPOS xmlns=\"http://www.modirum.com/schemas/vposxmlapi41\" xmlns:ns2=\"http://www.w3.org/2000/09/xmldsig#\"><Message version=\"2.1\" messageId=\"ADV12345678912345\" timeStamp=\"2026-03-29T01:44:29.071+02:00\"><Advice type=\"Sale\"><Authentication><Mid>9000000000</Mid></Authentication><OrderId>123465</OrderId><OrderAmount>446.4</OrderAmount><Currency>EUR</Currency><OrderTxId>12345678912345</OrderTxId><OrderTxStatus>ERROR</OrderTxStatus><PaymentTotal>446.4</PaymentTotal><TxId>12345678912345</TxId><TxStatus>ERROR</TxStatus><TxTotal>446.4</TxTotal><TxCurrency>EUR</TxCurrency></Advice></Message><Digest>K/OB9YUxjAABOJ+M8iFkHW755Dn/lfbep4aLF96RXKI=</Digest></VPOS>")

	resp, err := parseWebhookResponse(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OrderID != "123465" {
		t.Errorf("OrderID = %q, want 123465", resp.OrderID)
	}
	if resp.Status != cardlink.StatusError {
		t.Errorf("Status = %q, want ERROR", resp.Status)
	}
}
