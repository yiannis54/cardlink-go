package redirect

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/internal/digest"
)

func first(v url.Values, keys ...string) string {
	for _, k := range keys {
		if s := strings.TrimSpace(v.Get(k)); s != "" {
			return s
		}
	}
	return ""
}

// responsePreimage builds the digest preimage for standard redirect responses (Table 3).
func responsePreimage(v url.Values) string {
	return first(v, "version") +
		first(v, "mid") +
		first(v, "orderid") +
		first(v, "status") +
		first(v, "orderAmount") +
		first(v, "currency") +
		first(v, "paymentTotal") +
		first(v, "message") +
		first(v, "riskScore") +
		first(v, "payMethod") +
		first(v, "txId", "txid") +
		first(v, "paymentRef") +
		first(v, "extData")
}

// recurringNotificationPreimage builds the preimage for recurring notifications (Table 4).
func recurringNotificationPreimage(v url.Values) string {
	return first(v, "version") +
		first(v, "mid") +
		first(v, "orderid") +
		first(v, "status") +
		first(v, "orderAmount") +
		first(v, "currency") +
		first(v, "paymentTotal") +
		first(v, "message") +
		first(v, "riskScore") +
		first(v, "payMethod") +
		first(v, "txId", "txid") +
		first(v, "Sequence", "sequence") +
		first(v, "SeqTxId", "seqTxId", "seqtxid") +
		first(v, "paymentRef")
}

// VerifyResponse validates the digest on a standard payment callback.
func (s *Signer) VerifyResponse(v url.Values) (*Response, error) {
	if s.Config.SharedSecret == "" {
		return nil, fmt.Errorf("redirect: SharedSecret is required")
	}
	got := first(v, "digest")
	if got == "" {
		return nil, fmt.Errorf("redirect: missing digest")
	}
	pre := responsePreimage(v)
	expected := digest.Redirect(pre, s.Config.SharedSecret)
	if got != expected {
		return nil, fmt.Errorf("redirect: digest mismatch")
	}
	return parseResponse(v, got), nil
}

// VerifyRecurringNotification validates the digest on a scheduled recurring child notification.
// If version is missing or not "2", SHA-1 is used per Cardlink documentation.
func (s *Signer) VerifyRecurringNotification(v url.Values) (*RecurringNotification, error) {
	if s.Config.SharedSecret == "" {
		return nil, fmt.Errorf("redirect: SharedSecret is required")
	}
	got := first(v, "digest")
	if got == "" {
		return nil, fmt.Errorf("redirect: missing digest")
	}
	pre := recurringNotificationPreimage(v)
	ver := first(v, "version")
	var expected string
	if ver == "" || ver != "2" {
		expected = digest.RedirectSHA1(pre, s.Config.SharedSecret)
	} else {
		expected = digest.Redirect(pre, s.Config.SharedSecret)
	}
	if got != expected {
		return nil, fmt.Errorf("redirect: digest mismatch")
	}
	base := parseResponse(v, got)
	return &RecurringNotification{
		Response: *base,
		Sequence: first(v, "Sequence", "sequence"),
		SeqTxID:  first(v, "SeqTxId", "seqTxId", "seqtxid"),
	}, nil
}

func parseResponse(v url.Values, dig string) *Response {
	return &Response{
		Version:      first(v, "version"),
		MID:          first(v, "mid"),
		OrderID:      first(v, "orderid"),
		Status:       cardlink.ParseStatus(first(v, "status")),
		OrderAmount:  first(v, "orderAmount"),
		Currency:     first(v, "currency"),
		PaymentTotal: first(v, "paymentTotal"),
		Message:      first(v, "message"),
		RiskScore:    first(v, "riskScore"),
		PayMethod:    first(v, "payMethod"),
		TxID:         first(v, "txId", "txid"),
		PaymentRef:   first(v, "paymentRef"),
		ExtData:      first(v, "extData"),
		Digest:       dig,
	}
}
