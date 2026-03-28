package vposxml

import (
	"fmt"
	"time"

	"github.com/beevik/etree"
)

func newMessage(messageID, timeStamp string) *etree.Element {
	m := etree.NewElement("Message")
	m.CreateAttr("xmlns", vposNS)
	m.CreateAttr("xmlns:ns2", dsigNS)
	m.CreateAttr("messageId", messageID)
	m.CreateAttr("timeStamp", timeStamp)
	m.CreateAttr("version", "2.1")
	return m
}

// NewMessageID returns a message id in the form M{unixMilli} matching Cardlink samples.
func NewMessageID() string {
	return fmt.Sprintf("M%d", time.Now().UnixMilli())
}

// FormatTimeStamp formats t like Cardlink samples: 2023-04-20T16:59:44.802+03:00
func FormatTimeStamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.000-07:00")
}
