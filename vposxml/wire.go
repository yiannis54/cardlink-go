package vposxml

import (
	"bytes"
	"strings"

	"github.com/beevik/etree"
	"github.com/yiannis54/cardlink-go/internal/digest"
)

func wrapVPOS(message *etree.Element, digestVal string) (string, error) {
	doc := etree.NewDocument()
	doc.WriteSettings = etree.WriteSettings{CanonicalText: true}
	root := etree.NewElement("VPOS")
	root.CreateAttr("xmlns", vposNS)
	root.CreateAttr("xmlns:ns2", dsigNS)
	root.AddChild(message)
	dig := etree.NewElement("Digest")
	dig.SetText(digestVal)
	root.AddChild(dig)
	doc.SetRoot(root)
	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		return "", err
	}
	s := buf.String()
	if !strings.HasPrefix(s, "<?xml") {
		s = xmlHeader + s
	}
	return s, nil
}

const xmlHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"

func signMessage(message *etree.Element, secret string) (string, error) {
	c14n, err := digest.CanonicalXML10Message(message)
	if err != nil {
		return "", err
	}
	return digest.VPOS21(c14n, secret), nil
}
