package vposxml

import (
	"context"
	"time"

	"github.com/beevik/etree"
)

// RefundParams is a RefundRequest.
type RefundParams struct {
	MessageID   string
	TimeStamp   string
	MID         string
	OrderID     string
	OrderAmount string
	Currency    string
}

// RefundResponse is a parsed RefundResponse.
type RefundResponse struct {
	OrderID      string
	OrderAmount  string
	Currency     string
	PaymentTotal string
	Status       string
	TxID         string
	Description  string
	RawXML       []byte
}

// Refund executes a refund request.
func (c *Client) Refund(ctx context.Context, p RefundParams) (*RefundResponse, error) {
	if c.cfg.SharedSecret == "" {
		return nil, ErrMissingSecret
	}
	m, err := c.buildRefundMessage(p)
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
	return parseRefundResponse(raw)
}

func (c *Client) buildRefundMessage(p RefundParams) (*etree.Element, error) {
	mid := pickMID(c, p.MID)
	if mid == "" {
		return nil, ErrMissingMID
	}
	msgID := pickMsgID(p.MessageID)
	ts := pickTS(p.TimeStamp)
	m := newMessage(msgID, ts)
	rr := m.CreateElement("RefundRequest")
	rr.CreateElement("Authentication").CreateElement("Mid").SetText(mid)
	oi := rr.CreateElement("OrderInfo")
	oi.CreateElement("OrderId").SetText(p.OrderID)
	oi.CreateElement("OrderAmount").SetText(p.OrderAmount)
	oi.CreateElement("Currency").SetText(p.Currency)
	return m, nil
}

func parseRefundResponse(raw []byte) (*RefundResponse, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(raw); err != nil {
		return nil, err
	}
	root := doc.Root()
	var msg *etree.Element
	if root != nil && root.Tag == "VPOS" {
		msg = root.SelectElement("Message")
	} else if root != nil && root.Tag == "Message" {
		msg = root
	}
	if msg == nil {
		return &RefundResponse{RawXML: raw}, nil
	}
	if re := responseErrorFromMessage(msg); re != nil {
		return nil, re
	}
	rr := msg.SelectElement("RefundResponse")
	if rr == nil {
		return &RefundResponse{RawXML: raw}, nil
	}
	return &RefundResponse{
		OrderID:      textOrEmpty(rr, "./OrderId"),
		OrderAmount:  textOrEmpty(rr, "./OrderAmount"),
		Currency:     textOrEmpty(rr, "./Currency"),
		PaymentTotal: textOrEmpty(rr, "./PaymentTotal"),
		Status:       textOrEmpty(rr, "./Status"),
		TxID:         firstText(rr, "TxId"),
		Description:  textOrEmpty(rr, "./Description"),
		RawXML:       raw,
	}, nil
}

// CancelParams is a CancelRequest.
type CancelParams struct {
	MessageID   string
	TimeStamp   string
	MID         string
	OrderID     string
	OrderAmount string
	Currency    string
}

// CancelResponse is a parsed CancelResponse.
type CancelResponse struct {
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

// Cancel executes a cancel (void) request.
func (c *Client) Cancel(ctx context.Context, p CancelParams) (*CancelResponse, error) {
	if c.cfg.SharedSecret == "" {
		return nil, ErrMissingSecret
	}
	m, err := c.buildCancelMessage(p)
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
	return parseCancelResponse(raw)
}

func (c *Client) buildCancelMessage(p CancelParams) (*etree.Element, error) {
	mid := pickMID(c, p.MID)
	if mid == "" {
		return nil, ErrMissingMID
	}
	msgID := pickMsgID(p.MessageID)
	ts := pickTS(p.TimeStamp)
	m := newMessage(msgID, ts)
	cr := m.CreateElement("CancelRequest")
	cr.CreateElement("Authentication").CreateElement("Mid").SetText(mid)
	oi := cr.CreateElement("OrderInfo")
	oi.CreateElement("OrderId").SetText(p.OrderID)
	oi.CreateElement("OrderAmount").SetText(p.OrderAmount)
	oi.CreateElement("Currency").SetText(p.Currency)
	return m, nil
}

func parseCancelResponse(raw []byte) (*CancelResponse, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(raw); err != nil {
		return nil, err
	}
	root := doc.Root()
	var msg *etree.Element
	if root != nil && root.Tag == "VPOS" {
		msg = root.SelectElement("Message")
	} else if root != nil && root.Tag == "Message" {
		msg = root
	}
	if msg == nil {
		return &CancelResponse{RawXML: raw}, nil
	}
	if re := responseErrorFromMessage(msg); re != nil {
		return nil, re
	}
	cr := msg.SelectElement("CancelResponse")
	if cr == nil {
		return &CancelResponse{RawXML: raw}, nil
	}
	return &CancelResponse{
		OrderID:      textOrEmpty(cr, "./OrderId"),
		OrderAmount:  textOrEmpty(cr, "./OrderAmount"),
		Currency:     textOrEmpty(cr, "./Currency"),
		PaymentTotal: textOrEmpty(cr, "./PaymentTotal"),
		Status:       textOrEmpty(cr, "./Status"),
		TxID:         firstText(cr, "TxId"),
		PaymentRef:   textOrEmpty(cr, "./PaymentRef"),
		Description:  textOrEmpty(cr, "./Description"),
		RawXML:       raw,
	}, nil
}

func pickMID(c *Client, mid string) string {
	if mid != "" {
		return mid
	}
	return c.cfg.MID
}

func pickMsgID(id string) string {
	if id != "" {
		return id
	}
	return NewMessageID()
}

func pickTS(ts string) string {
	if ts != "" {
		return ts
	}
	return FormatTimeStamp(time.Now())
}
