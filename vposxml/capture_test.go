package vposxml

import (
	"testing"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/internal/digest"
)

func TestBuildCaptureMessage_CanonicalMatchesGolden(t *testing.T) {
	cfg := cardlink.Config{
		MID:          "0101119349",
		SharedSecret: "ignored",
		Environment:  cardlink.Sandbox,
		Partner:      cardlink.Cardlink,
	}
	c := NewClient(cfg)
	m, err := c.buildCaptureMessage(CaptureParams{
		MessageID:   "M1681999184802",
		TimeStamp:   "2023-04-20T16:59:44.802+03:00",
		MID:         "0101119349",
		OrderID:     "O230419155347",
		OrderAmount: "50.0",
		Currency:    "EUR",
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := digest.CanonicalXML10Message(m)
	if err != nil {
		t.Fatal(err)
	}
	// No whitespace-only text between elements — otherwise C14N includes those text nodes.
	const goldenCompact = `<Message xmlns="http://www.modirum.com/schemas/vposxmlapi41" xmlns:ns2="http://www.w3.org/2000/09/xmldsig#" messageId="M1681999184802" timeStamp="2023-04-20T16:59:44.802+03:00" version="2.1"><CaptureRequest><Authentication><Mid>0101119349</Mid></Authentication><OrderInfo><OrderId>O230419155347</OrderId><OrderAmount>50.0</OrderAmount><Currency>EUR</Currency></OrderInfo></CaptureRequest></Message>`
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

func TestSignMessage_VPOS21Stable(t *testing.T) {
	cfg := cardlink.Config{MID: "0101119349", SharedSecret: "s3cret"}
	c := NewClient(cfg)
	m, err := c.buildCaptureMessage(CaptureParams{
		MessageID:   "M1681999184802",
		TimeStamp:   "2023-04-20T16:59:44.802+03:00",
		OrderID:     "O230419155347",
		OrderAmount: "50.0",
		Currency:    "EUR",
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
		t.Fatalf("digest not stable")
	}
}
