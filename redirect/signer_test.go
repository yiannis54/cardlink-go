package redirect

import (
	"net/url"
	"testing"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/internal/digest"
)

func TestSign_GoldenSaleExample(t *testing.T) {
	cfg := cardlink.Config{
		MID:          "0101119349",
		SharedSecret: "Cardlink1",
		Environment:  cardlink.Sandbox,
		Partner:      cardlink.Worldline,
	}
	s := NewSigner(cfg)
	req := &PaymentRequest{
		Version:        "2",
		Lang:           "en",
		DeviceCategory: "0",
		OrderID:        "O170911143656",
		OrderDesc:      "Test order some items",
		OrderAmount:    "0.12",
		Currency:       "EUR",
		PayerEmail:     "cardlink@cardlink.gr",
		PayerPhone:     "30-6900000000",
		BillCountry:    "GR",
		BillZip:        "12345",
		BillCity:       "Athens",
		BillAddress:    "Street 45",
		ConfirmURL:     "https://ecommerce-test.cardlink.gr/vpostestsv4/shops/shopdemo.jsp?cmd=confirm",
		CancelURL:      "https://ecommerce-test.cardlink.gr/vpostestsv4/shops/shopdemo.jsp?cmd=cancel",
	}
	fields, err := s.Sign(req)
	if err != nil {
		t.Fatal(err)
	}
	want := "ybXX2tQkFlxzHM5SjH0oGrD9zms21SUQnwkYaFrnGdc="
	if fields["digest"] != want {
		t.Fatalf("digest = %q, want %q", fields["digest"], want)
	}
}

func TestVerifyResponse_Golden(t *testing.T) {
	cfg := cardlink.Config{SharedSecret: "Cardlink1"}
	s := NewSigner(cfg)
	v := url.Values{}
	v.Set("version", "2")
	v.Set("mid", "0101119349")
	v.Set("orderid", "O170911143656")
	v.Set("status", "CAPTURED")
	v.Set("orderAmount", "0.12")
	v.Set("currency", "EUR")
	v.Set("paymentTotal", "0.12")
	v.Set("message", "OK, 00 - Approved")
	v.Set("riskScore", "0")
	v.Set("payMethod", "visa")
	v.Set("txId", "926012471")
	v.Set("paymentRef", "138104")
	v.Set("digest", "FpwgGyCRwhmF6CWtRFLqfkuQpdPyX8Xh3tJg3E891SA=")

	resp, err := s.VerifyResponse(v)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != cardlink.StatusCaptured {
		t.Fatalf("status %v", resp.Status)
	}
}

func TestVerifyRecurringNotification_SHA256(t *testing.T) {
	cfg := cardlink.Config{SharedSecret: "secret"}
	s := NewSigner(cfg)
	v := url.Values{}
	v.Set("version", "2")
	v.Set("mid", "1")
	v.Set("orderid", "O1")
	v.Set("status", "CAPTURED")
	v.Set("orderAmount", "1.00")
	v.Set("currency", "EUR")
	v.Set("paymentTotal", "1.00")
	v.Set("txId", "9")
	v.Set("Sequence", "2")
	v.Set("SeqTxId", "99")
	v.Set("paymentRef", "1")
	pre := recurringNotificationPreimage(v)
	v.Set("digest", digest.Redirect(pre, cfg.SharedSecret))

	_, err := s.VerifyRecurringNotification(v)
	if err != nil {
		t.Fatal(err)
	}
}

func TestVerifyRecurringNotification_SHA1Fallback(t *testing.T) {
	cfg := cardlink.Config{SharedSecret: "secret"}
	s := NewSigner(cfg)
	v := url.Values{}
	// no version
	v.Set("mid", "1")
	v.Set("orderid", "O1")
	v.Set("status", "CAPTURED")
	v.Set("orderAmount", "1.00")
	v.Set("currency", "EUR")
	v.Set("paymentTotal", "1.00")
	v.Set("txId", "9")
	v.Set("Sequence", "2")
	v.Set("SeqTxId", "99")
	v.Set("paymentRef", "1")
	pre := recurringNotificationPreimage(v)
	v.Set("digest", digest.RedirectSHA1(pre, cfg.SharedSecret))

	_, err := s.VerifyRecurringNotification(v)
	if err != nil {
		t.Fatal(err)
	}
}
