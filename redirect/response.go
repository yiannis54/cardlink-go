package redirect

import "github.com/yiannis54/cardlink-go/cardlink"

// Response is the standard redirect POST callback (confirm or cancel URL).
type Response struct {
	Version      string
	MID          string
	OrderID      string
	Status       cardlink.Status
	OrderAmount  string
	Currency     string
	PaymentTotal string
	Message      string
	RiskScore    string
	PayMethod    string
	TxID         string
	PaymentRef   string
	ExtData      string
	Digest       string
}

// RecurringNotification is the recurring child notification POST (recurring notify URL).
type RecurringNotification struct {
	Response
	Sequence string
	SeqTxID  string
}
