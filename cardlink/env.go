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

// ParseEnvironment parses a string into an Environment.
func ParseEnvironment(s string) Environment {
	switch s {
	case "sandbox":
		return Sandbox
	case "production":
		return Production
	default:
		return Sandbox
	}
}
