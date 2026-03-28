package cardlink

import (
	"fmt"
	"regexp"
	"strings"
)

var orderIDRe = regexp.MustCompile(`^[A-Za-z0-9]+$`)

// ValidateOrderID checks order id per redirect rules: alphanumeric, max 50 (45 if recurring).
func ValidateOrderID(id string, recurring bool) error {
	max := 50
	if recurring {
		max = 45
	}
	if id == "" {
		return fmt.Errorf("cardlink: empty orderid")
	}
	if len(id) > max {
		return fmt.Errorf("cardlink: orderid exceeds %d characters", max)
	}
	if !orderIDRe.MatchString(id) {
		return fmt.Errorf("cardlink: orderid must contain only letters and numbers")
	}
	return nil
}

// FormatOrderAmount returns amount string without thousands separators, using dot as decimal (e.g. "10346.78").
func FormatOrderAmount(amount string) string {
	return strings.TrimSpace(strings.ReplaceAll(amount, ",", ""))
}
