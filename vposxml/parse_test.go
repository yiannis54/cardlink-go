package vposxml

import (
	"errors"
	"testing"

	"github.com/yiannis54/cardlink-go/cardlink"
)

func TestParseCaptureResponse_ErrorEnvelope(t *testing.T) {
	raw := []byte(`<VPOS xmlns="http://www.modirum.com/schemas/vposxmlapi41"><Message version="1.0" messageId="M1" timeStamp="2024-01-01T00:00:00+02:00"><ErrorCode>XE</ErrorCode><ErrorMessage>invalid XML</ErrorMessage><Description>detail</Description><OriginalXML>&lt;x/&gt;</OriginalXML></Message></VPOS>`)
	_, err := parseCaptureResponse(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	var re *cardlink.ResponseError
	if !errors.As(err, &re) {
		t.Fatalf("expected ResponseError, got %v", err)
	}
	if re.ErrorCode != "XE" || re.ErrorMessage != "invalid XML" || re.Description != "detail" {
		t.Fatalf("unexpected fields: %+v", re)
	}
	if re.OriginalXML == "" {
		t.Fatal("expected OriginalXML for version 1.0 envelope")
	}
}

func TestResponseErrorFromMessage_NoErrorMessage(t *testing.T) {
	if responseErrorFromMessage(nil) != nil {
		t.Fatal("expected nil")
	}
}
