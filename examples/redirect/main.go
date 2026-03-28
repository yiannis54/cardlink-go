// Command redirect demonstrates signing a Cardlink redirect payment request.
package main

import (
	"fmt"
	"log"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/redirect"
)

func main() {
	cfg := cardlink.Config{
		MID:          "0101119349",
		SharedSecret: "replace-with-shared-secret",
		Environment:  cardlink.Sandbox,
		Partner:      cardlink.Worldline,
	}

	url, err := redirect.NewSigner(cfg).
		CheckoutURL(&redirect.PaymentRequest{
			OrderID:     "ORDER123",
			OrderDesc:   "Demo order",
			OrderAmount: "1.00",
			Currency:    "EUR",
			PayerEmail:  "buyer@example.com",
			ConfirmURL:  "https://merchant.example.com/pay/ok",
			CancelURL:   "https://merchant.example.com/pay/cancel",
		})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(url)
}
