// Example: handle Cardlink confirmUrl — browser return and Modirum VPOS background POST.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/redirect"
)

func main() {
	cfg := cardlink.Config{
		MID:          "your-mid",
		SharedSecret: "your-shared-secret",
		Environment:  cardlink.Sandbox,
		Partner:      cardlink.Cardlink,
	}
	signer := redirect.NewSigner(cfg)

	http.HandleFunc("/pay/confirm", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		v := r.Form

		resp, err := signer.VerifyResponse(v)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Idempotent fulfillment keyed by resp.OrderID (and handle unknown order → 406).
		_ = resp

		if redirect.IsServerToServerCallback(r) {
			log.Println("background Modirum VPOS notification")
		} else {
			log.Println("browser return to confirmUrl")
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
