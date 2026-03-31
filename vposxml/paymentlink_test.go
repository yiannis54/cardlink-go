package vposxml

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/internal/digest"
)

func boolPtr(v bool) *bool { return &v }

func TestBuildPaymentLinkMessage_BasicGolden(t *testing.T) {
	cfg := cardlink.Config{MID: "90006064", SharedSecret: "ignored"}
	c := NewClient(cfg)
	m, err := c.buildPaymentLinkMessage(PaymentLinkParams{
		MessageID:           "M1728981255715",
		TimeStamp:           "2024-10-15T11:34:15.715+03:00",
		OrderID:             "1728981205998",
		OrderDesc:           "",
		OrderAmount:         "3.0",
		Currency:            "EUR",
		PayerEmail:          "test@example.com",
		TxType:              PaymentLinkTxPayment,
		LinkValidityDays:    5,
		MailLinkIfValidMail: boolPtr(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := digest.CanonicalXML10Message(m)
	if err != nil {
		t.Fatal(err)
	}
	const goldenCompact = `<Message xmlns="http://www.modirum.com/schemas/vposxmlapi41" xmlns:ns2="http://www.w3.org/2000/09/xmldsig#" messageId="M1728981255715" timeStamp="2024-10-15T11:34:15.715+03:00" version="2.1"><PaymentLinkRequest><Authentication><Mid>90006064</Mid></Authentication><OrderInfo><OrderId>1728981205998</OrderId><OrderDesc></OrderDesc><OrderAmount>3.0</OrderAmount><Currency>EUR</Currency><PayerEmail>test@example.com</PayerEmail></OrderInfo><PaymentInfo></PaymentInfo><TxType>PAYMENT</TxType><LinkValidityDays>5</LinkValidityDays><MailLinkIfValidMail>true</MailLinkIfValidMail></PaymentLinkRequest></Message>`
	gc, err := digest.ParseMessageElement(goldenCompact)
	if err != nil {
		t.Fatal(err)
	}
	want, err := digest.CanonicalXML10Message(gc)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != string(want) {
		t.Errorf("canonical mismatch\n got: %s\nwant: %s", out, want)
	}
}

func TestBuildPaymentLinkMessage_Installment(t *testing.T) {
	cfg := cardlink.Config{MID: "0101118297", SharedSecret: "ignored"}
	c := NewClient(cfg)
	m, err := c.buildPaymentLinkMessage(PaymentLinkParams{
		MessageID:         "M1751528174293",
		TimeStamp:         "2025-07-03T10:36:14.293+03:00",
		OrderID:           "1751528090098",
		OrderDesc:         "",
		OrderAmount:       "4.0",
		Currency:          "EUR",
		PayerEmail:        "test@example.com",
		TxType:            PaymentLinkTxPayment,
		LinkValidityDays:  7,
		InstallmentOffset: 0,
		InstallmentPeriod: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := digest.CanonicalXML10Message(m)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "<InstallmentParameters>") {
		t.Error("expected <InstallmentParameters>")
	}
	if !strings.Contains(s, "<ExtInstallmentoffset>0</ExtInstallmentoffset>") {
		t.Error("expected ExtInstallmentoffset=0")
	}
	if !strings.Contains(s, "<ExtInstallmentperiod>2</ExtInstallmentperiod>") {
		t.Error("expected ExtInstallmentperiod=2")
	}
	if strings.Contains(s, "<MailLinkIfValidMail>") {
		t.Error("MailLinkIfValidMail should be omitted when nil")
	}
}

func TestBuildPaymentLinkMessage_WithBillingAddress(t *testing.T) {
	cfg := cardlink.Config{MID: "90006064", SharedSecret: "ignored"}
	c := NewClient(cfg)
	m, err := c.buildPaymentLinkMessage(PaymentLinkParams{
		MessageID:   "M100",
		TimeStamp:   "2024-10-15T11:34:15.715+03:00",
		OrderID:     "O1",
		OrderAmount: "1.0",
		Currency:    "EUR",
		PayerEmail:  "test@example.com",
		TxType:      PaymentLinkTxPayment,
		BillingAddress: &Address{
			Country: "GR",
			State:   "",
			Zip:     "12345",
			City:    "Athens",
			Street:  "Test 12",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := digest.CanonicalXML10Message(m)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "<BillingAddress>") {
		t.Error("expected <BillingAddress>")
	}
	if !strings.Contains(s, "<country>GR</country>") {
		t.Error("expected <country>GR</country>")
	}
	if !strings.Contains(s, "<address>Test 12</address>") {
		t.Error("expected <address>Test 12</address>")
	}
}

func TestBuildPaymentLinkMessage_WithLang(t *testing.T) {
	cfg := cardlink.Config{MID: "90006064", SharedSecret: "ignored"}
	c := NewClient(cfg)
	m, err := c.buildPaymentLinkMessage(PaymentLinkParams{
		MessageID:   "M200",
		TimeStamp:   "2024-10-15T11:34:15.715+03:00",
		Lang:        "el",
		OrderID:     "O2",
		OrderAmount: "1.0",
		Currency:    "EUR",
		PayerEmail:  "test@example.com",
		TxType:      PaymentLinkTxPayment,
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := digest.CanonicalXML10Message(m)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `lang="el"`) {
		t.Error("expected lang attribute in canonical output")
	}
}

func TestParsePaymentLinkResponse_Success(t *testing.T) {
	raw := []byte(`<VPOS xmlns="http://www.modirum.com/schemas/vposxmlapi41" xmlns:ns2="http://www.w3.org/2000/09/xmldsig#"><Message version="2.1" messageId="M1728981255715" timeStamp="2024-10-15T11:34:15.975+03:00"><PaymentLinkResponse><OrderId>1728981205998</OrderId><OrderAmount>3.0</OrderAmount><Currency>EUR</Currency><PaymentTotal>3.0</PaymentTotal><Status>EXECWAIT</Status><TxId>92639555657341</TxId><Description>OK, link created mailed</Description><PaymentLink>https://eurocommerce-test.cardlink.gr/vpos/Paylink/b033573194386268462f9034e0391f36</PaymentLink><LinkMailed>true</LinkMailed></PaymentLinkResponse></Message><Digest>IwQfgRMoiKvklLjzFBD+GphxTJDewH1beTcxNf+no6Y=</Digest></VPOS>`)
	resp, err := parsePaymentLinkResponse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.OrderID != "1728981205998" {
		t.Errorf("OrderID = %q", resp.OrderID)
	}
	if resp.Status != cardlink.StatusExecWait {
		t.Errorf("Status = %q, want EXECWAIT", resp.Status)
	}
	if resp.TxID != "92639555657341" {
		t.Errorf("TxID = %q", resp.TxID)
	}
	if resp.PaymentLink != "https://eurocommerce-test.cardlink.gr/vpos/Paylink/b033573194386268462f9034e0391f36" {
		t.Errorf("PaymentLink = %q", resp.PaymentLink)
	}
	if !resp.LinkMailed {
		t.Error("expected LinkMailed=true")
	}
	if resp.Description != "OK, link created mailed" {
		t.Errorf("Description = %q", resp.Description)
	}
}

func TestParsePaymentLinkResponse_LinkNotMailed(t *testing.T) {
	raw := []byte(`<VPOS xmlns="http://www.modirum.com/schemas/vposxmlapi41"><Message version="2.1" messageId="M1" timeStamp="2025-07-03T10:36:14.599+03:00"><PaymentLinkResponse><OrderId>O1</OrderId><OrderAmount>4.0</OrderAmount><Currency>EUR</Currency><PaymentTotal>4.0</PaymentTotal><Status>EXECWAIT</Status><TxId>123</TxId><Description>OK, link created</Description><PaymentLink>https://example.com/pay/abc</PaymentLink><LinkMailed>false</LinkMailed></PaymentLinkResponse></Message><Digest>x</Digest></VPOS>`)
	resp, err := parsePaymentLinkResponse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.LinkMailed {
		t.Error("expected LinkMailed=false")
	}
}

func TestParsePaymentLinkResponse_ErrorEnvelope(t *testing.T) {
	raw := []byte(`<VPOS xmlns="http://www.modirum.com/schemas/vposxmlapi41"><Message version="1.0" messageId="M1" timeStamp="2024-01-01T00:00:00+02:00"><ErrorCode>XE</ErrorCode><ErrorMessage>invalid XML</ErrorMessage><Description>detail</Description><OriginalXML>&lt;x/&gt;</OriginalXML></Message></VPOS>`)
	_, err := parsePaymentLinkResponse(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	var re *cardlink.ResponseError
	if !errors.As(err, &re) {
		t.Fatalf("expected ResponseError, got %v", err)
	}
	if re.ErrorCode != "XE" {
		t.Errorf("ErrorCode = %q", re.ErrorCode)
	}
}

func TestPaymentLink_MissingSecret(t *testing.T) {
	cfg := cardlink.Config{MID: "123", SharedSecret: ""}
	c := NewClient(cfg)
	_, err := c.PaymentLink(context.Background(), PaymentLinkParams{TxType: PaymentLinkTxPayment})
	if !errors.Is(err, ErrMissingSecret) {
		t.Fatalf("expected ErrMissingSecret, got %v", err)
	}
}

func TestPaymentLink_MissingTxType(t *testing.T) {
	cfg := cardlink.Config{MID: "123", SharedSecret: "s3cret"}
	c := NewClient(cfg)
	_, err := c.PaymentLink(context.Background(), PaymentLinkParams{
		OrderID:     "O1",
		OrderAmount: "1.0",
		Currency:    "EUR",
		PayerEmail:  "test@example.com",
	})
	if !errors.Is(err, ErrMissingTxType) {
		t.Fatalf("expected errMissingTxType, got %v", err)
	}
}

func TestSignPaymentLinkMessage_Stable(t *testing.T) {
	cfg := cardlink.Config{MID: "90006064", SharedSecret: "s3cret"}
	c := NewClient(cfg)
	m, err := c.buildPaymentLinkMessage(PaymentLinkParams{
		MessageID:        "M1728981255715",
		TimeStamp:        "2024-10-15T11:34:15.715+03:00",
		OrderID:          "1728981205998",
		OrderAmount:      "3.0",
		Currency:         "EUR",
		PayerEmail:       "test@example.com",
		TxType:           PaymentLinkTxPayment,
		LinkValidityDays: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	d1, err := signMessage(m, "s3cret")
	if err != nil {
		t.Fatal(err)
	}
	d2, err := signMessage(m, "s3cret")
	if err != nil {
		t.Fatal(err)
	}
	if d1 != d2 {
		t.Fatalf("digest not stable: %s != %s", d1, d2)
	}
}
