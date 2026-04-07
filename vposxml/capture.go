package vposxml

import (
	"context"
	"time"

	"github.com/beevik/etree"
)

// CaptureParams is a VPOS XML CaptureRequest.
type CaptureParams struct {
	MessageID   string
	TimeStamp   string
	MID         string
	OrderID     string
	OrderAmount string
	Currency    string
	Var1        string
	Var2        string
}

// CaptureResponse is parsed CaptureResponse.
// When the gateway returns an error envelope (ErrorMessage), the operation returns a non-nil error wrapping ResponseError.
type CaptureResponse struct {
	OrderID      string
	OrderAmount  string
	Currency     string
	PaymentTotal string
	Status       string
	TxID         string
	PaymentRef   string
	Description  string
	RawXML       []byte
}

func (c *Client) buildCaptureMessage(p CaptureParams) (*etree.Element, error) {
	mid := p.MID
	if mid == "" {
		mid = c.cfg.MID
	}
	if mid == "" {
		return nil, ErrMissingMID
	}
	msgID := p.MessageID
	if msgID == "" {
		msgID = NewMessageID()
	}
	ts := p.TimeStamp
	if ts == "" {
		ts = FormatTimeStamp(time.Now())
	}
	m := newMessage(msgID, ts)
	cr := m.CreateElement("CaptureRequest")
	auth := cr.CreateElement("Authentication")
	auth.CreateElement("Mid").SetText(mid)
	oi := cr.CreateElement("OrderInfo")
	oi.CreateElement("OrderId").SetText(p.OrderID)
	oi.CreateElement("OrderAmount").SetText(p.OrderAmount)
	oi.CreateElement("Currency").SetText(p.Currency)
	if p.Var1 != "" {
		oi.CreateElement("Var1").SetText(p.Var1)
	}
	if p.Var2 != "" {
		oi.CreateElement("Var2").SetText(p.Var2)
	}
	return m, nil
}

// Capture executes a capture request.
func (c *Client) Capture(ctx context.Context, p CaptureParams) (*CaptureResponse, error) {
	if c.cfg.SharedSecret == "" {
		return nil, ErrMissingSecret
	}
	m, err := c.buildCaptureMessage(p)
	if err != nil {
		return nil, err
	}
	dig, err := signMessage(m, c.cfg.SharedSecret)
	if err != nil {
		return nil, err
	}
	xml, err := wrapVPOS(m, dig)
	if err != nil {
		return nil, err
	}
	raw, err := c.responseBytes(ctx, xml, true)
	if err != nil {
		return nil, err
	}
	return parseCaptureResponse(raw)
}
