// Command vposxml demonstrates building a VPOS XML capture request body (no HTTP by default).
package main

import (
	"fmt"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/vposxml"
)

func main() {
	cfg := cardlink.Config{
		MID:          "0101119349",
		SharedSecret: "replace-with-shared-secret",
		Environment:  cardlink.Sandbox,
		Partner:      cardlink.Cardlink,
	}
	_ = vposxml.NewClient(cfg)
	fmt.Println("Configure MID and SharedSecret, then use vposxml.Client.Capture / Status / IRISSale / RecurringOperation / PaymentLink.")
}
