# cardlink-go

Go SDK for [Cardlink](https://developer.cardlink.gr/) **redirect** (HTML form POST + digest) and **VPOS XML 2.1** (HTTP XML + Canonical XML 1.0 digest), with sandbox/production endpoints for Cardlink, Nexi, and Worldline.

## Install

```bash
go get github.com/yiannis54/cardlink-go@latest
```

## Redirect payments

1. Configure [`cardlink.Config`](cardlink/config.go) with `MID`, `SharedSecret`, `Environment` (`Sandbox` / `Production`), and `BusinessPartner` (`Cardlink`, `Nexi`, `Worldline`).
2. Build a [`redirect.PaymentRequest`](redirect/payment.go), then [`redirect.Signer.Sign`](redirect/signer.go) to obtain POST fields including `digest` (server-side only).
3. POST the browser to `cfg.RedirectURL()` (path `/vpos/shophandlermpi`), or render [`redirect.FormHTML`](redirect/html.go).
4. On return to `confirmUrl` / `cancelUrl`, call [`Signer.VerifyResponse`](redirect/verify.go). Use [`IsServerToServerCallback`](redirect/callback.go) to detect delayed `Modirum VPOS` background posts. Cardlink may retry server-to-server posts: respond with **200** for recognized orders, **406** if unknown, **400** for malformed input (per gateway documentation); implement idempotency on `orderid`.

**Subscriptions (master):** set `ExtRecurringFrequency`, `ExtRecurringEndDate` on `PaymentRequest`. For **unscheduled** recurring, set `Var6` to `rcauto=false` per Cardlink docs. **Installments** use `ExtInstallmentOffset` / `ExtInstallmentPeriod` — not combinable with recurring.

**IRIS (redirect):** set `PayMethod` to `IRIS` where supported. **IRIS QR** generation for display uses VPOS XML — see below.

## VPOS XML 2.1

Use [`vposxml.NewClient`](vposxml/client.go) with the same `cardlink.Config`. The client POSTs to `cfg.VPOSXMLURL()` (path `/vpos/xmlpayvpos`).

Supported operations:

- [`Client.Capture`](vposxml/capture.go)
- [`Client.Status`](vposxml/status.go)
- [`Client.Refund`](vposxml/refund_cancel.go) / [`Client.Cancel`](vposxml/refund_cancel.go)
- [`Client.IRISSale`](vposxml/iris.go) (`PayMethod=iris`, `PaymentOption=irisQr`)
- [`Client.RecurringOperation`](vposxml/recurring.go) (`RecurringChild`, `Cancel`)

Digest verification uses **inclusive Canonical XML 1.0** on the `<Message>` element and `Base64(SHA256(c14n || sharedSecret))` (see Cardlink direct integration docs).

**Out of scope for this module:** VPOS XML 4.1 (XML-DSig), redirect tokenization tables, 3DS/MPI direct card entry, Payment Link XML, CSE, mass files.

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

## Endpoints

Sandbox hosts (examples):

- Cardlink: `https://ecommerce-test.cardlink.gr`
- Nexi: `https://alphaecommerce-test.cardlink.gr`
- Worldline: `https://eurocommerce-test.cardlink.gr`

Production hosts omit `-test` where applicable — see [`cardlink/config.go`](cardlink/config.go). Override with `RedirectBaseURL` / `VPOSXMLBaseURL` if needed.

## Examples

See [`examples/`](examples/):

- `examples/redirect` — sign a payment and print hidden fields
- `examples/vposxml` — build a capture request (no network unless credentials set)
- `examples/webhook` — handle Cardlink confirmUrl

## Documentation

- [Redirect integration](https://developer.cardlink.gr/api_products_categories/redirect-integration/)
- [Recurring transactions](https://developer.cardlink.gr/documentation_categories/recurring-transactions/)
- [VPOS XML sample requests](https://developer.cardlink.gr/api_products_categories/vpos-xml-requests/)

## License

Use at your own risk; no warranty. Validate all flows against Cardlink test credentials.
