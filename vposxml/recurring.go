package vposxml

import (
	"context"

	"github.com/beevik/etree"
	"github.com/yiannis54/cardlink-go/cardlink"
)

// RecurringOperation is the operation name for RecurringOperationRequest.
type RecurringOperation string

const (
	OpRecurringChild RecurringOperation = "RecurringChild"
	OpCancel         RecurringOperation = "Cancel"
)

// RecurringOperationParams is a RecurringOperationRequest.
type RecurringOperationParams struct {
	MessageID string
	TimeStamp string
	MID       string
	OrderID   string // master recurring order id
	Operation RecurringOperation
	// Optional: set when charging a different amount than the master (unscheduled recurring child).
	ChildOrderDesc   string
	ChildOrderAmount string
	ChildCurrency    string
}

// RecurringOperationResponse is a parsed RecurringOperationResponse.
type RecurringOperationResponse struct {
	OrderID      string
	OrderAmount  string
	Currency     string
	PaymentTotal string
	Status       cardlink.Status
	TxID         string
	Sequence     string
	PaymentRef   string
	RiskScore    string
	Description  string
	RawXML       []byte
}

// RecurringOperation executes RecurringChild or Cancel.
func (c *Client) RecurringOperation(ctx context.Context, p RecurringOperationParams) (*RecurringOperationResponse, error) {
	if c.cfg.SharedSecret == "" {
		return nil, ErrMissingSecret
	}
	m, err := c.buildRecurringMessage(p)
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
	return parseRecurringOperationResponse(raw)
}

func parseRecurringOperationResponse(raw []byte) (*RecurringOperationResponse, error) {
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
		return &RecurringOperationResponse{RawXML: raw}, nil
	}
	if re := responseErrorFromMessage(msg); re != nil {
		return nil, re
	}
	rr := msg.SelectElement("RecurringOperationResponse")
	if rr == nil {
		return &RecurringOperationResponse{RawXML: raw}, nil
	}
	return &RecurringOperationResponse{
		OrderID:      textOrEmpty(rr, "./OrderId"),
		OrderAmount:  textOrEmpty(rr, "./OrderAmount"),
		Currency:     textOrEmpty(rr, "./Currency"),
		PaymentTotal: textOrEmpty(rr, "./PaymentTotal"),
		Status:       cardlink.ParseStatus(textOrEmpty(rr, "./Status")),
		TxID:         firstText(rr, "TxId"),
		Sequence:     textOrEmpty(rr, "./Sequence"),
		PaymentRef:   textOrEmpty(rr, "./PaymentRef"),
		RiskScore:    textOrEmpty(rr, "./RiskScore"),
		Description:  textOrEmpty(rr, "./Description"),
		RawXML:       raw,
	}, nil
}

func (c *Client) buildRecurringMessage(p RecurringOperationParams) (*etree.Element, error) {
	mid := pickMID(c, p.MID)
	if mid == "" {
		return nil, ErrMissingMID
	}
	if p.Operation == "" {
		return nil, errMissingOperation
	}
	msgID := pickMsgID(p.MessageID)
	ts := pickTS(p.TimeStamp)
	m := newMessage(msgID, ts)
	rr := m.CreateElement("RecurringOperationRequest")
	rr.CreateElement("Authentication").CreateElement("Mid").SetText(mid)
	ti := rr.CreateElement("TransactionInfo")
	ti.CreateElement("OrderId").SetText(p.OrderID)
	rr.CreateElement("Operation").SetText(string(p.Operation))
	if p.ChildOrderAmount != "" || p.ChildCurrency != "" {
		oi := rr.CreateElement("OrderInfo")
		oi.CreateElement("OrderId").SetText(p.OrderID)
		if p.ChildOrderDesc != "" {
			oi.CreateElement("OrderDesc").SetText(p.ChildOrderDesc)
		}
		if p.ChildOrderAmount != "" {
			oi.CreateElement("OrderAmount").SetText(p.ChildOrderAmount)
		}
		if p.ChildCurrency != "" {
			oi.CreateElement("Currency").SetText(p.ChildCurrency)
		}
	}
	return m, nil
}
