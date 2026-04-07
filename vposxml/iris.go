package vposxml

import (
	"context"
	"strings"

	"github.com/beevik/etree"
	"github.com/yiannis54/cardlink-go/cardlink"
)

// IRISSaleParams is a SaleRequest for IRIS QR (Worldline Greece).
type IRISSaleParams struct {
	MessageID   string
	TimeStamp   string
	MID         string
	OrderID     string
	OrderDesc   string
	OrderAmount string
	Currency    string
	Var1        string
	Var2        string
}

// IRISSaleResponse is a parsed SaleResponse for IRIS.
type IRISSaleResponse struct {
	OrderID      string
	OrderAmount  string
	Currency     string
	PaymentTotal string
	Status       cardlink.Status
	TxID         string
	IRISQR       string // data URI or raw payload from Attribute IRIS-QR
	IRISTXID     string
	Description  string
	RawXML       []byte
}

// IRISSale executes an IRIS QR sale request.
func (c *Client) IRISSale(ctx context.Context, p IRISSaleParams) (*IRISSaleResponse, error) {
	if c.cfg.SharedSecret == "" {
		return nil, ErrMissingSecret
	}
	m, err := c.buildIRISSaleMessage(p)
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
	return parseIRISSaleResponse(raw)
}

func (c *Client) buildIRISSaleMessage(p IRISSaleParams) (*etree.Element, error) {
	mid := pickMID(c, p.MID)
	if mid == "" {
		return nil, ErrMissingMID
	}
	msgID := pickMsgID(p.MessageID)
	ts := pickTS(p.TimeStamp)
	m := newMessage(msgID, ts)
	sr := m.CreateElement("SaleRequest")
	sr.CreateElement("Authentication").CreateElement("Mid").SetText(mid)
	oi := sr.CreateElement("OrderInfo")
	oi.CreateElement("OrderId").SetText(p.OrderID)
	oi.CreateElement("OrderDesc").SetText(p.OrderDesc)
	oi.CreateElement("OrderAmount").SetText(p.OrderAmount)
	oi.CreateElement("Currency").SetText(p.Currency)
	if p.Var1 != "" {
		oi.CreateElement("Var1").SetText(p.Var1)
	}
	if p.Var2 != "" {
		oi.CreateElement("Var2").SetText(p.Var2)
	}
	pi := sr.CreateElement("PaymentInfo")
	pi.CreateElement("PayMethod").SetText("iris")
	pi.CreateElement("PaymentOption").SetText("irisQr")
	return m, nil
}

func parseIRISSaleResponse(raw []byte) (*IRISSaleResponse, error) {
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
		return &IRISSaleResponse{RawXML: raw}, nil
	}
	if re := responseErrorFromMessage(msg); re != nil {
		return nil, re
	}
	sr := msg.SelectElement("SaleResponse")
	if sr == nil {
		return &IRISSaleResponse{RawXML: raw}, nil
	}
	irisQR := ""
	irisTx := ""
	for _, el := range sr.ChildElements() {
		if el.Tag != "Attribute" {
			continue
		}
		switch el.SelectAttrValue("name", "") {
		case "IRIS-QR":
			irisQR = strings.TrimSpace(el.Text())
		case "IRIS-TXID":
			irisTx = strings.TrimSpace(el.Text())
		}
	}
	return &IRISSaleResponse{
		OrderID:      textOrEmpty(sr, "./OrderId"),
		OrderAmount:  textOrEmpty(sr, "./OrderAmount"),
		Currency:     textOrEmpty(sr, "./Currency"),
		PaymentTotal: textOrEmpty(sr, "./PaymentTotal"),
		Status:       cardlink.ParseStatus(textOrEmpty(sr, "./Status")),
		TxID:         firstText(sr, "TxId"),
		IRISQR:       irisQR,
		IRISTXID:     irisTx,
		Description:  textOrEmpty(sr, "./Description"),
		RawXML:       raw,
	}, nil
}
