package cardlink

// Status is a payment or transaction status from gateway responses.
type Status string

const (
	StatusAuthorized         Status = "AUTHORIZED"
	StatusCaptured           Status = "CAPTURED"
	StatusCanceled           Status = "CANCELED"
	StatusRefused            Status = "REFUSED"
	StatusRefusedRisk        Status = "REFUSEDRISK"
	StatusError              Status = "ERROR"
	StatusProcessing         Status = "PROCESSING"
	StatusExecWait           Status = "EXECWAIT"
	StatusPreprocess         Status = "PREPROCESS"
	StatusPreprocessTimedOut Status = "PREPROCESS-TIMEDOUT"
	StatusInWallet           Status = "INWALLET"
	StatusExecWaitTimedOut   Status = "EXECWAIT-TIMEDOUT"
)

// ParseStatus returns s as Status if known; otherwise returns Status(s) for forward compatibility.
func ParseStatus(s string) Status {
	switch Status(s) {
	case StatusAuthorized, StatusCaptured, StatusCanceled, StatusRefused, StatusRefusedRisk,
		StatusError, StatusProcessing, StatusExecWait, StatusPreprocess, StatusPreprocessTimedOut,
		StatusInWallet, StatusExecWaitTimedOut:
		return Status(s)
	default:
		return Status(s)
	}
}

func (s Status) String() string {
	return string(s)
}
