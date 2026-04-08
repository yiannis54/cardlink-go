package vposxml

import (
	"context"
	"net/http"
)

// ClientAPI defines the consumer-facing operations exposed by Client.
// Depend on this interface in application code to allow easy mocking.
type ClientAPI interface {
	Capture(ctx context.Context, p CaptureParams) (*CaptureResponse, error)
	Status(ctx context.Context, p StatusParams) (*StatusResponse, error)
	Refund(ctx context.Context, p RefundParams) (*RefundResponse, error)
	Cancel(ctx context.Context, p CancelParams) (*CancelResponse, error)
	PaymentLink(ctx context.Context, p PaymentLinkParams) (*PaymentLinkResponse, error)
	IRISSale(ctx context.Context, p IRISSaleParams) (*IRISSaleResponse, error)
	RecurringOperation(ctx context.Context, p RecurringOperationParams) (*RecurringOperationResponse, error)
	VerifyWebhook(xmlBody []byte) (*WebhookResponse, error)
	VerifyWebhookRequest(r *http.Request) (*WebhookResponse, error)
	RawPOST(ctx context.Context, vposXML string, verifyResponse bool) ([]byte, error)
}

var _ ClientAPI = (*Client)(nil)
