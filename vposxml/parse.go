package vposxml

import (
	"strings"

	"github.com/beevik/etree"

	"github.com/yiannis54/cardlink-go/cardlink"
)

func textOrEmpty(el *etree.Element, path string) string {
	if el == nil {
		return ""
	}
	if x := el.FindElement(path); x != nil {
		return strings.TrimSpace(x.Text())
	}
	return ""
}

// responseErrorFromMessage returns a typed gateway error when the Message contains ErrorMessage.
func responseErrorFromMessage(msg *etree.Element) *cardlink.ResponseError {
	if msg == nil {
		return nil
	}
	em := msg.SelectElement("ErrorMessage")
	if em == nil {
		return nil
	}
	re := &cardlink.ResponseError{
		ErrorCode:    textOrEmpty(msg, "./ErrorCode"),
		ErrorMessage: strings.TrimSpace(em.Text()),
		Description:  textOrEmpty(msg, "./Description"),
	}
	if msg.SelectAttrValue("version", "") == "1.0" {
		re.OriginalXML = textOrEmpty(msg, "./OriginalXML")
	}
	return re
}

func parseCaptureResponse(raw []byte) (*CaptureResponse, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(raw); err != nil {
		return nil, err
	}
	root := doc.Root()
	if root == nil {
		return &CaptureResponse{RawXML: raw}, nil
	}
	var msg *etree.Element
	if root.Tag == "VPOS" {
		msg = root.SelectElement("Message")
	} else {
		msg = root
	}
	if msg == nil {
		return &CaptureResponse{RawXML: raw}, nil
	}
	if re := responseErrorFromMessage(msg); re != nil {
		return nil, re
	}
	cr := msg.SelectElement("CaptureResponse")
	if cr == nil {
		return &CaptureResponse{RawXML: raw}, nil
	}
	return &CaptureResponse{
		OrderID:      textOrEmpty(cr, "./OrderId"),
		OrderAmount:  textOrEmpty(cr, "./OrderAmount"),
		Currency:     textOrEmpty(cr, "./Currency"),
		PaymentTotal: textOrEmpty(cr, "./PaymentTotal"),
		Status:       textOrEmpty(cr, "./Status"),
		TxID:         firstText(cr, "TxId", "TxID"),
		PaymentRef:   textOrEmpty(cr, "./PaymentRef"),
		Description:  textOrEmpty(cr, "./Description"),
		RawXML:       raw,
	}, nil
}

func firstText(el *etree.Element, names ...string) string {
	for _, n := range names {
		if x := el.SelectElement(n); x != nil {
			return strings.TrimSpace(x.Text())
		}
	}
	return ""
}

// varText looks for a free variable first as a direct child element (e.g. <Var1>),
// then as an <Attribute name="VAR1"> element (the format used in StatusResponse and webhooks).
func varText(el *etree.Element, elemName, attrName string) string {
	if v := textOrEmpty(el, "./"+elemName); v != "" {
		return v
	}
	for _, child := range el.ChildElements() {
		if child.Tag == "Attribute" && child.SelectAttrValue("name", "") == attrName {
			return strings.TrimSpace(child.Text())
		}
	}
	return ""
}
