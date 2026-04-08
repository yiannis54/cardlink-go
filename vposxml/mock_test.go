package vposxml

import (
	"context"
	"strings"
	"testing"
)

func TestMockClientAPI_CaptureDispatch(t *testing.T) {
	t.Parallel()

	mock := &MockClientAPI{
		CaptureFunc: func(_ context.Context, p CaptureParams) (*CaptureResponse, error) {
			return &CaptureResponse{
				OrderID: p.OrderID,
				Status:  "CAPTURED",
			}, nil
		},
	}

	resp, err := mock.Capture(context.Background(), CaptureParams{OrderID: "ORD-1"})
	if err != nil {
		t.Fatalf("Capture() unexpected err: %v", err)
	}
	if resp == nil {
		t.Fatal("Capture() returned nil response")
	}
	if resp.OrderID != "ORD-1" {
		t.Fatalf("Capture().OrderID = %q, want %q", resp.OrderID, "ORD-1")
	}
	if resp.Status != "CAPTURED" {
		t.Fatalf("Capture().Status = %q, want %q", resp.Status, "CAPTURED")
	}
}

func TestMockClientAPI_UnsetMethodReturnsError(t *testing.T) {
	t.Parallel()

	mock := &MockClientAPI{}
	_, err := mock.Status(context.Background(), StatusParams{})
	if err == nil {
		t.Fatal("Status() expected error for unset function")
	}
	if !strings.Contains(err.Error(), "Status is not implemented") {
		t.Fatalf("Status() error = %q, want contains %q", err.Error(), "Status is not implemented")
	}
}

func TestMockClientAPI_RawPOSTDispatch(t *testing.T) {
	t.Parallel()

	mock := &MockClientAPI{
		RawPOSTFunc: func(_ context.Context, vposXML string, verifyResponse bool) ([]byte, error) {
			if !verifyResponse {
				t.Fatal("RawPOST verifyResponse should be true")
			}
			return []byte(vposXML), nil
		},
	}

	raw, err := mock.RawPOST(context.Background(), "<VPOS/>", true)
	if err != nil {
		t.Fatalf("RawPOST() unexpected err: %v", err)
	}
	if got, want := string(raw), "<VPOS/>"; got != want {
		t.Fatalf("RawPOST() = %q, want %q", got, want)
	}
}
