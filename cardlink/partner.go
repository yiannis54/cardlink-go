package cardlink

// BusinessPartner selects Cardlink, Nexi, or Worldline gateway hosts.
type BusinessPartner int

const (
	Cardlink BusinessPartner = iota
	Nexi
	Worldline
)

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
