# cardlink-go

Go SDK for [Cardlink](https://developer.cardlink.gr/) **VPOS XML 2.1** (HTTP XML + Canonical XML 1.0 digest), with sandbox/production endpoints for Cardlink, Nexi, and Worldline.

## Install

```bash
go get github.com/yiannis54/cardlink-go@latest
```

## VPOS XML 2.1

Use [`vposxml.NewClient`](vposxml/client.go) with the same `cardlink.Config`. The client POSTs to `cfg.VPOSXMLURL()` (path `/vpos/xmlpayvpos`).

Supported operations:

- [`Client.Capture`](vposxml/capture.go)
- [`Client.Status`](vposxml/status.go)
- [`Client.Refund`](vposxml/refund_cancel.go) / [`Client.Cancel`](vposxml/refund_cancel.go)
- [`Client.IRISSale`](vposxml/iris.go) (`PayMethod=iris`, `PaymentOption=irisQr`)
- [`Client.RecurringOperation`](vposxml/recurring.go) (`RecurringChild`, `Cancel`)
- [`Client.PaymentLink`](vposxml/paymentlink.go) (create and optionally email a payment link)

Digest verification uses **inclusive Canonical XML 1.0** on the `<Message>` element and `Base64(SHA256(c14n || sharedSecret))` (see Cardlink direct integration docs).

**Out of scope for this module:** VPOS XML 4.1 (XML-DSig) and other direct-integration features not implemented (e.g. tokenizer/token tables, 3DS/MPI direct card entry, CSE, mass files).

### Quick start: IRIS payment

IRIS QR sale uses [`Client.IRISSale`](vposxml/iris.go). Use sandbox credentials from Cardlink and set `Partner` to `Worldline` for the Worldline Greece host. Replace placeholders with your MID, shared secret, and order fields.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/vposxml"
)

func main() {
	cfg := cardlink.Config{
		MID:          "your-mid",
		SharedSecret: "your-shared-secret",
		Environment:  cardlink.Sandbox,
		Partner:      cardlink.Worldline,
	}
	client := vposxml.NewClient(cfg)

	resp, err := client.IRISSale(context.Background(), vposxml.IRISSaleParams{
		OrderID:     "order-001",
		OrderDesc:   "IRIS test",
		OrderAmount: "1.00",
		Currency:    "978", // EUR (ISO 4217 numeric)
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("status:", resp.Status)
	fmt.Println("IRIS-QR:", resp.IRISQR)
}
```

### Quick start: Payment Link

Use [`Client.PaymentLink`](vposxml/paymentlink.go) to create a payment link and optionally have the gateway email it to the payer. Supports installments and billing address.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/vposxml"
)

func main() {
	cfg := cardlink.Config{
		MID:          "your-mid",
		SharedSecret: "your-shared-secret",
		Environment:  cardlink.Sandbox,
		Partner:      cardlink.Worldline,
	}
	client := vposxml.NewClient(cfg)

	mailLink := true
	resp, err := client.PaymentLink(context.Background(), vposxml.PaymentLinkParams{
		OrderID:             "order-002",
		OrderAmount:         "25.00",
		Currency:            "EUR",
		PayerEmail:          "customer@example.com",
		TxType:              vposxml.PaymentLinkTxPayment,
		LinkValidityDays:    5,
		MailLinkIfValidMail: &mailLink,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("status:", resp.Status)
	fmt.Println("link:", resp.PaymentLink)
	fmt.Println("mailed:", resp.LinkMailed)
}
```

## Endpoints

Sandbox hosts (examples):

- Cardlink: `https://ecommerce-test.cardlink.gr`
- Nexi: `https://alphaecommerce-test.cardlink.gr`
- Worldline: `https://eurocommerce-test.cardlink.gr`

Production hosts omit `-test` where applicable — see [`cardlink/config.go`](cardlink/config.go). Override with `VPOSXMLBaseURL` if needed.

## Examples

See [`examples/`](examples/):

- `examples/vposxml` — build a capture request (no network unless credentials set)
- `examples/webhook` — verify VPOS XML callbacks/notifications

## Documentation

- [Recurring transactions](https://developer.cardlink.gr/documentation_categories/recurring-transactions/)
- [VPOS XML sample requests](https://developer.cardlink.gr/api_products_categories/vpos-xml-requests/)
- [Payment Link through XML API](https://developer.cardlink.gr/api_products_categories/payment-link/)

## License

Use at your own risk; no warranty. Validate all flows against Cardlink test credentials.
