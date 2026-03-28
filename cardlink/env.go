package cardlink

// Environment selects sandbox (test) vs production hosts.
type Environment int

const (
	Sandbox Environment = iota
	Production
)

func (e Environment) String() string {
	switch e {
	case Sandbox:
		return "sandbox"
	case Production:
		return "production"
	default:
		return "unknown"
	}
}
