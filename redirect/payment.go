package redirect

// PaymentRequest contains all fields that participate in the redirect digest, in table order.
// Optional fields left empty contribute "" to the digest preimage.
type PaymentRequest struct {
	Version string // required, typically "2"
	// MID is set from cardlink.Config.MID during Sign if empty
	MID string

	Lang                  string
	DeviceCategory        string
	OrderID               string
	OrderDesc             string
	OrderAmount           string
	Currency              string
	PayerEmail            string
	PayerPhone            string
	BillCountry           string
	BillState             string
	BillZip               string
	BillCity              string
	BillAddress           string
	Weight                string
	Dimensions            string
	ShipCountry           string
	ShipState             string
	ShipZip               string
	ShipCity              string
	ShipAddress           string
	AddFraudScore         string
	MaxPayRetries         string
	Reject3dsU            string
	PayMethod             string
	TrType                string
	ExtInstallmentOffset  string
	ExtInstallmentPeriod  string
	ExtRecurringFrequency string
	ExtRecurringEndDate   string
	BlockScore            string
	CssURL                string
	ConfirmURL            string
	CancelURL             string
	Var1, Var2, Var3      string
	Var4, Var5, Var6      string
	Var7, Var8, Var9      string
}
