package main

import (
	"context"
	"errors"
	"testing"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/vposxml"
)

type captureService struct {
	api vposxml.ClientAPI
}

func (s captureService) CreatePaymentLink(ctx context.Context, orderID string) (string, error) {
	resp, err := s.api.PaymentLink(ctx, vposxml.PaymentLinkParams{
		OrderID:     orderID,
		OrderAmount: "10.00",
		Currency:    "EUR",
		PayerEmail:  "buyer@example.com",
		TxType:      vposxml.PaymentLinkTxPayment,
	})
	if err != nil {
		return "", err
	}
	return resp.PaymentLink, nil
}

func TestCaptureService_WithMock_Success(t *testing.T) {
	t.Parallel()

	svc := captureService{
		api: &vposxml.MockClientAPI{
			PaymentLinkFunc: func(_ context.Context, p vposxml.PaymentLinkParams) (*vposxml.PaymentLinkResponse, error) {
				return &vposxml.PaymentLinkResponse{
					OrderID:     p.OrderID,
					PaymentLink: "https://pay.example/ORD-200",
					Status:      cardlink.StatusExecWait,
				}, nil
			},
		},
	}

	link, err := svc.CreatePaymentLink(context.Background(), "ORD-200")
	if err != nil {
		t.Fatalf("CreatePaymentLink() unexpected err: %v", err)
	}
	if got, want := link, "https://pay.example/ORD-200"; got != want {
		t.Fatalf("CreatePaymentLink() = %q, want %q", got, want)
	}
}

func TestCaptureService_WithMock_Failure(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("gateway unavailable")
	svc := captureService{
		api: &vposxml.MockClientAPI{
			PaymentLinkFunc: func(_ context.Context, _ vposxml.PaymentLinkParams) (*vposxml.PaymentLinkResponse, error) {
				return nil, wantErr
			},
		},
	}

	_, err := svc.CreatePaymentLink(context.Background(), "ORD-201")
	if !errors.Is(err, wantErr) {
		t.Fatalf("CreatePaymentLink() error = %v, want %v", err, wantErr)
	}
}
