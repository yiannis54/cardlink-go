package digest

import (
	"fmt"

	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
)

// CanonicalXML10Message canonicalizes the <Message> element using inclusive Canonical XML 1.0
// (http://www.w3.org/TR/2001/REC-xml-c14n-20010315) as required for VPOS XML 2.1 digest input.
func CanonicalXML10Message(message *etree.Element) ([]byte, error) {
	if message == nil || message.Tag != "Message" {
		return nil, fmt.Errorf("digest: root element must be Message")
	}
	c := dsig.MakeC14N10RecCanonicalizer()
	return c.Canonicalize(message)
}

// ParseMessageElement parses XML that contains a single root <Message> (optionally with XML declaration).
func ParseMessageElement(xmlStr string) (*etree.Element, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlStr); err != nil {
		return nil, err
	}
	root := doc.Root()
	if root == nil || root.Tag != "Message" {
		return nil, fmt.Errorf("digest: expected Message root")
	}
	return root, nil
}
