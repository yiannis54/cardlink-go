package cardlink

import "fmt"

// BusinessPartner selects Cardlink, Nexi, or Worldline gateway hosts.
type BusinessPartner int

const (
	Cardlink BusinessPartner = iota
	Nexi
	Worldline
)

// ParseBusinessPartner parses a string into a BusinessPartner.
func ParseBusinessPartner(s string) (BusinessPartner, error) {
	switch s {
	case "cardlink":
		return Cardlink, nil
	case "nexi":
		return Nexi, nil
	case "worldline":
		return Worldline, nil
	default:
		return 0, fmt.Errorf("unknown business partner: %s", s)
	}
}

// String returns the string representation of a BusinessPartner.
func (p BusinessPartner) String() string {
	switch p {
	case Cardlink:
		return "cardlink"
	case Nexi:
		return "nexi"
	case Worldline:
		return "worldline"
	default:
		return "unknown"
	}
}
