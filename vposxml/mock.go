package vposxml

import (
	"context"
	"fmt"
	"net/http"
)

// MockClientAPI is a hand-written mock for ClientAPI.
// Set only the function fields you need in each test.
type MockClientAPI struct {
	CaptureFunc            func(ctx context.Context, p CaptureParams) (*CaptureResponse, error)
	StatusFunc             func(ctx context.Context, p StatusParams) (*StatusResponse, error)
	RefundFunc             func(ctx context.Context, p RefundParams) (*RefundResponse, error)
	CancelFunc             func(ctx context.Context, p CancelParams) (*CancelResponse, error)
	PaymentLinkFunc        func(ctx context.Context, p PaymentLinkParams) (*PaymentLinkResponse, error)
	IRISSaleFunc           func(ctx context.Context, p IRISSaleParams) (*IRISSaleResponse, error)
	RecurringOperationFunc func(ctx context.Context, p RecurringOperationParams) (*RecurringOperationResponse, error)
	VerifyWebhookFunc      func(xmlBody []byte) (*WebhookResponse, error)
	VerifyWebhookReqFunc   func(r *http.Request) (*WebhookResponse, error)
	RawPOSTFunc            func(ctx context.Context, vposXML string, verifyResponse bool) ([]byte, error)
}

var _ ClientAPI = (*MockClientAPI)(nil)

func mockNotImplemented(method string) error {
	return fmt.Errorf("vposxml mock: %s is not implemented", method)
}

func (m *MockClientAPI) Capture(ctx context.Context, p CaptureParams) (*CaptureResponse, error) {
	if m.CaptureFunc == nil {
		return nil, mockNotImplemented("Capture")
	}
	return m.CaptureFunc(ctx, p)
}

func (m *MockClientAPI) Status(ctx context.Context, p StatusParams) (*StatusResponse, error) {
	if m.StatusFunc == nil {
		return nil, mockNotImplemented("Status")
	}
	return m.StatusFunc(ctx, p)
}

func (m *MockClientAPI) Refund(ctx context.Context, p RefundParams) (*RefundResponse, error) {
	if m.RefundFunc == nil {
		return nil, mockNotImplemented("Refund")
	}
	return m.RefundFunc(ctx, p)
}

func (m *MockClientAPI) Cancel(ctx context.Context, p CancelParams) (*CancelResponse, error) {
	if m.CancelFunc == nil {
		return nil, mockNotImplemented("Cancel")
	}
	return m.CancelFunc(ctx, p)
}

func (m *MockClientAPI) PaymentLink(ctx context.Context, p PaymentLinkParams) (*PaymentLinkResponse, error) {
	if m.PaymentLinkFunc == nil {
		return nil, mockNotImplemented("PaymentLink")
	}
	return m.PaymentLinkFunc(ctx, p)
}

func (m *MockClientAPI) IRISSale(ctx context.Context, p IRISSaleParams) (*IRISSaleResponse, error) {
	if m.IRISSaleFunc == nil {
		return nil, mockNotImplemented("IRISSale")
	}
	return m.IRISSaleFunc(ctx, p)
}

func (m *MockClientAPI) RecurringOperation(ctx context.Context, p RecurringOperationParams) (*RecurringOperationResponse, error) {
	if m.RecurringOperationFunc == nil {
		return nil, mockNotImplemented("RecurringOperation")
	}
	return m.RecurringOperationFunc(ctx, p)
}

func (m *MockClientAPI) VerifyWebhook(xmlBody []byte) (*WebhookResponse, error) {
	if m.VerifyWebhookFunc == nil {
		return nil, mockNotImplemented("VerifyWebhook")
	}
	return m.VerifyWebhookFunc(xmlBody)
}

func (m *MockClientAPI) VerifyWebhookRequest(r *http.Request) (*WebhookResponse, error) {
	if m.VerifyWebhookReqFunc == nil {
		return nil, mockNotImplemented("VerifyWebhookRequest")
	}
	return m.VerifyWebhookReqFunc(r)
}

func (m *MockClientAPI) RawPOST(ctx context.Context, vposXML string, verifyResponse bool) ([]byte, error) {
	if m.RawPOSTFunc == nil {
		return nil, mockNotImplemented("RawPOST")
	}
	return m.RawPOSTFunc(ctx, vposXML, verifyResponse)
}
