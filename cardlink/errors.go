package cardlink

import "fmt"

// Error codes from Cardlink VPOS XML / gateway documentation.
const (
	ErrM1 = "M1" // invalid MID
	ErrM2 = "M2" // authentication failed / wrong digest
	ErrSE = "SE" // system error
	ErrXE = "XE" // invalid XML
	ErrI0 = "I0" // unsupported request
	ErrI1 = "I1" // invalid or missing data
	ErrI2 = "I2" // invalid installments
	ErrI3 = "I3" // invalid recurring parameters
	ErrI4 = "I4" // invalid card data
	ErrI5 = "I5" // invalid expiration
	ErrI6 = "I6" // payment method mismatch
	ErrO1 = "O1" // operation not allowed
	ErrO2 = "O2" // original transaction not found
)

// ResponseError is returned when the gateway signals an application-level error (e.g. XML error envelope).
// Inspect with errors.As into a *ResponseError when handling vposxml.Client method errors.
type ResponseError struct {
	ErrorCode    string
	ErrorMessage string
	Description  string
	OriginalXML  string
}

func (e *ResponseError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("cardlink: %s %s: %s", e.ErrorCode, e.ErrorMessage, e.Description)
}
