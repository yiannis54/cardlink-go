// Example: handle Cardlink confirmUrl — browser return and Modirum VPOS background POST.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/vposxml"
)

func main() {
	cfg := cardlink.Config{
		MID:          "your-mid",
		SharedSecret: "your-shared-secret",
		Environment:  cardlink.Sandbox,
		Partner:      cardlink.Cardlink,
	}
	client := vposxml.NewClient(cfg)

	http.HandleFunc("/pay/confirm", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		resp, err := client.VerifyWebhookRequest(r)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Idempotent fulfillment keyed by resp.OrderID (and handle unknown order → 406).
		_ = resp

		log.Printf("vposxml callback: order=%s status=%s", resp.OrderID, resp.Status)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
