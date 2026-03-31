package vposxml

import (
	"context"
	"strconv"
	"strings"

	"github.com/beevik/etree"

	"github.com/yiannis54/cardlink-go/cardlink"
)

// PaymentLinkTxType is the transaction type for a PaymentLinkRequest.
type PaymentLinkTxType string

const (
	PaymentLinkTxPayment PaymentLinkTxType = "PAYMENT"
	PaymentLinkTxPreauth PaymentLinkTxType = "PAYMENT_PREAUTH"
)

// Address holds billing or shipping address fields.
// Element names match the Cardlink XML schema (lowercase).
type Address struct {
	Country string
	State   string
	Zip     string
	City    string
	Street  string // emitted as <address>
}

// PaymentLinkParams is a VPOS XML PaymentLinkRequest.
type PaymentLinkParams struct {
	MessageID string
	TimeStamp string
	MID       string
	Lang      string // optional ISO 639-1 language code for email template (e.g. "el", "en")

	OrderID     string
	OrderDesc   string
	OrderAmount string
	Currency    string
	PayerEmail  string
	PayerPhone  string

	BillingAddress *Address

	TxType              PaymentLinkTxType
	LinkValidityDays    int
	MailLinkIfValidMail *bool

	InstallmentOffset int
	InstallmentPeriod int

	PayMethod     string
	PaymentOption string

	Var1 string
	Var2 string
}

// PaymentLinkResponse is a parsed PaymentLinkResponse.
type PaymentLinkResponse struct {
	OrderID      string
	OrderAmount  string
	Currency     string
	PaymentTotal string
	Status       cardlink.Status
	TxID         string
	Description  string
	PaymentLink  string
	LinkMailed   bool
	ErrorCode    string
	RawXML       []byte
}

// PaymentLink creates and optionally emails a payment link via VPOS XML.
func (c *Client) PaymentLink(ctx context.Context, p PaymentLinkParams) (*PaymentLinkResponse, error) {
	if c.cfg.SharedSecret == "" {
		return nil, ErrMissingSecret
	}
	m, err := c.buildPaymentLinkMessage(p)
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
	return parsePaymentLinkResponse(raw)
}

func (c *Client) buildPaymentLinkMessage(p PaymentLinkParams) (*etree.Element, error) {
	mid := pickMID(c, p.MID)
	if mid == "" {
		return nil, ErrMissingMID
	}
	if p.TxType == "" {
		return nil, ErrMissingTxType
	}
	msgID := pickMsgID(p.MessageID)
	ts := pickTS(p.TimeStamp)
	m := newMessage(msgID, ts)
	if p.Lang != "" {
		m.CreateAttr("lang", p.Lang)
	}

	plr := m.CreateElement("PaymentLinkRequest")

	plr.CreateElement("Authentication").CreateElement("Mid").SetText(mid)

	oi := plr.CreateElement("OrderInfo")
	oi.CreateElement("OrderId").SetText(p.OrderID)
	oi.CreateElement("OrderDesc").SetText(p.OrderDesc)
	oi.CreateElement("OrderAmount").SetText(p.OrderAmount)
	oi.CreateElement("Currency").SetText(p.Currency)
	oi.CreateElement("PayerEmail").SetText(p.PayerEmail)
	if p.PayerPhone != "" {
		oi.CreateElement("PayerPhone").SetText(p.PayerPhone)
	}
	if p.BillingAddress != nil {
		ba := oi.CreateElement("BillingAddress")
		ba.CreateElement("country").SetText(p.BillingAddress.Country)
		ba.CreateElement("state").SetText(p.BillingAddress.State)
		ba.CreateElement("zip").SetText(p.BillingAddress.Zip)
		ba.CreateElement("city").SetText(p.BillingAddress.City)
		ba.CreateElement("address").SetText(p.BillingAddress.Street)
	}

	pi := plr.CreateElement("PaymentInfo")
	if p.InstallmentPeriod > 0 {
		ip := pi.CreateElement("InstallmentParameters")
		ip.CreateElement("ExtInstallmentoffset").SetText(strconv.Itoa(p.InstallmentOffset))
		ip.CreateElement("ExtInstallmentperiod").SetText(strconv.Itoa(p.InstallmentPeriod))
	}

	plr.CreateElement("TxType").SetText(string(p.TxType))

	if p.LinkValidityDays > 0 {
		plr.CreateElement("LinkValidityDays").SetText(strconv.Itoa(p.LinkValidityDays))
	}
	if p.MailLinkIfValidMail != nil {
		v := "false"
		if *p.MailLinkIfValidMail {
			v = "true"
		}
		plr.CreateElement("MailLinkIfValidMail").SetText(v)
	}

	return m, nil
}

func parsePaymentLinkResponse(raw []byte) (*PaymentLinkResponse, error) {
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
		return &PaymentLinkResponse{RawXML: raw}, nil
	}
	if re := responseErrorFromMessage(msg); re != nil {
		return nil, re
	}
	plr := msg.SelectElement("PaymentLinkResponse")
	if plr == nil {
		return &PaymentLinkResponse{RawXML: raw}, nil
	}
	return &PaymentLinkResponse{
		OrderID:      textOrEmpty(plr, "./OrderId"),
		OrderAmount:  textOrEmpty(plr, "./OrderAmount"),
		Currency:     textOrEmpty(plr, "./Currency"),
		PaymentTotal: textOrEmpty(plr, "./PaymentTotal"),
		Status:       cardlink.ParseStatus(textOrEmpty(plr, "./Status")),
		TxID:         firstText(plr, "TxId"),
		Description:  textOrEmpty(plr, "./Description"),
		PaymentLink:  textOrEmpty(plr, "./PaymentLink"),
		LinkMailed:   strings.EqualFold(textOrEmpty(plr, "./LinkMailed"), "true"),
		ErrorCode:    textOrEmpty(plr, "./ErrorCode"),
		RawXML:       raw,
	}, nil
}
