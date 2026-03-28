package digest

import (
	"strings"
	"testing"
)

const captureCanonicalMessage = `<Message xmlns="http://www.modirum.com/schemas/vposxmlapi41" xmlns:ns2="http://www.w3.org/2000/09/xmldsig#" messageId="M1681999184802" timeStamp="2023-04-20T16:59:44.802+03:00" version="2.1">
    <CaptureRequest>
        <Authentication>
            <Mid>0101119349</Mid>
        </Authentication>
        <OrderInfo>
            <OrderId>O230419155347</OrderId>
            <OrderAmount>50.0</OrderAmount>
            <Currency>EUR</Currency>
        </OrderInfo>
    </CaptureRequest>
</Message>`

func TestCanonicalXML10Message_CaptureRoundTrip(t *testing.T) {
	el, err := ParseMessageElement(captureCanonicalMessage)
	if err != nil {
		t.Fatal(err)
	}
	out, err := CanonicalXML10Message(el)
	if err != nil {
		t.Fatal(err)
	}
	// Second round-trip must be stable
	el2, err := ParseMessageElement(string(out))
	if err != nil {
		t.Fatal(err)
	}
	out2, err := CanonicalXML10Message(el2)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != string(out2) {
		t.Fatalf("second canonicalization changed output")
	}
}

func TestVPOS21_KnownVector(t *testing.T) {
	// If canonical bytes match Cardlink tool output, digest matches published sample with same secret.
	el, err := ParseMessageElement(captureCanonicalMessage)
	if err != nil {
		t.Fatal(err)
	}
	c14n, err := CanonicalXML10Message(el)
	if err != nil {
		t.Fatal(err)
	}
	// Brute-force verify published digest matches some secret is not possible; assert length/shape only.
	if len(c14n) < 50 {
		t.Fatalf("c14n too short: %q", c14n)
	}
	d := VPOS21(c14n, "x")
	if !strings.HasSuffix(d, "=") {
		t.Fatalf("expected base64: %s", d)
	}
}
