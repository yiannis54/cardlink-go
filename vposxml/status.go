package vposxml

import (
	"context"
	"time"

	"github.com/beevik/etree"
)

// StatusParams is a StatusRequest.
type StatusParams struct {
	MessageID string
	TimeStamp string
	MID       string
	OrderID   string
	Var1      string
	Var2      string
}

// StatusResponse is a parsed StatusResponse.
type StatusResponse struct {
	OrderID       string
	OrderAmount   string
	Currency      string
	PaymentTotal  string
	Status        string
	TxID          string
	PaymentRef    string
	RiskScore     string
	Description   string
	PaymentMethod string
	RawXML        []byte
}

func (c *Client) buildStatusMessage(p StatusParams) (*etree.Element, error) {
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
	sr := m.CreateElement("StatusRequest")
	sr.CreateElement("Authentication").CreateElement("Mid").SetText(mid)
	ti := sr.CreateElement("TransactionInfo")
	ti.CreateElement("OrderId").SetText(p.OrderID)
	if p.Var1 != "" {
		ti.CreateElement("Var1").SetText(p.Var1)
	}
	if p.Var2 != "" {
		ti.CreateElement("Var2").SetText(p.Var2)
	}
	return m, nil
}

// Status executes a status request.
func (c *Client) Status(ctx context.Context, p StatusParams) (*StatusResponse, error) {
	if c.cfg.SharedSecret == "" {
		return nil, ErrMissingSecret
	}
	m, err := c.buildStatusMessage(p)
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
	return parseStatusResponse(raw)
}

func parseStatusResponse(raw []byte) (*StatusResponse, error) {
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
		return &StatusResponse{RawXML: raw}, nil
	}
	if re := responseErrorFromMessage(msg); re != nil {
		return nil, re
	}
	sr := msg.SelectElement("StatusResponse")
	if sr == nil {
		return &StatusResponse{RawXML: raw}, nil
	}
	td := sr.SelectElement("TransactionDetails")
	if td == nil {
		return &StatusResponse{RawXML: raw}, nil
	}
	return &StatusResponse{
		OrderID:       textOrEmpty(td, "./OrderId"),
		OrderAmount:   textOrEmpty(td, "./OrderAmount"),
		Currency:      textOrEmpty(td, "./Currency"),
		PaymentTotal:  textOrEmpty(td, "./PaymentTotal"),
		Status:        textOrEmpty(td, "./Status"),
		TxID:          firstText(td, "TxId"),
		PaymentRef:    textOrEmpty(td, "./PaymentRef"),
		RiskScore:     textOrEmpty(td, "./RiskScore"),
		Description:   textOrEmpty(td, "./Description"),
		PaymentMethod: textOrEmpty(td, "./PaymentMethod"),
		RawXML:        raw,
	}, nil
}
