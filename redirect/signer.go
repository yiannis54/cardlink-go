// Package redirect implements Cardlink redirect integration: payment initiation field ordering,
// SHA-256 digest signing, response verification (including recurring notifications with SHA-1 fallback),
// and optional HTML form rendering.
package redirect

import (
	"fmt"

	"github.com/yiannis54/cardlink-go/cardlink"
	"github.com/yiannis54/cardlink-go/internal/digest"
)

// Signer signs redirect payment requests using Config.MID and Config.SharedSecret.
type Signer struct {
	Config cardlink.Config
}

// NewSigner returns a Signer for the given merchant configuration.
func NewSigner(cfg cardlink.Config) *Signer {
	return &Signer{Config: cfg}
}

// Sign validates the request, computes the digest, and returns all POST fields including digest.
func (s *Signer) Sign(req *PaymentRequest) (map[string]string, error) {
	if s.Config.SharedSecret == "" {
		return nil, fmt.Errorf("redirect: SharedSecret is required")
	}
	mid := req.MID
	if mid == "" {
		mid = s.Config.MID
	}
	if mid == "" {
		return nil, fmt.Errorf("redirect: MID is required (set PaymentRequest.MID or Config.MID)")
	}
	version := req.Version
	if version == "" {
		version = "2"
	}
	if err := cardlink.ValidateOrderID(req.OrderID, hasRecurring(req)); err != nil {
		return nil, err
	}
	if req.OrderDesc == "" || req.OrderAmount == "" || req.Currency == "" || req.PayerEmail == "" {
		return nil, fmt.Errorf("redirect: orderDesc, orderAmount, currency, and payerEmail are required")
	}
	if req.ConfirmURL == "" || req.CancelURL == "" {
		return nil, fmt.Errorf("redirect: confirmUrl and cancelUrl are required")
	}
	if err := mutualExclusiveInstallmentRecurring(req); err != nil {
		return nil, err
	}

	reqCopy := *req
	reqCopy.Version = version
	pre := redirectPreimage(mid, &reqCopy)
	dig := digest.Redirect(pre, s.Config.SharedSecret)

	fields := map[string]string{
		"version":               version,
		"mid":                   mid,
		"lang":                  req.Lang,
		"deviceCategory":        req.DeviceCategory,
		"orderid":               req.OrderID,
		"orderDesc":             req.OrderDesc,
		"orderAmount":           cardlink.FormatOrderAmount(req.OrderAmount),
		"currency":              req.Currency,
		"payerEmail":            req.PayerEmail,
		"payerPhone":            req.PayerPhone,
		"billCountry":           req.BillCountry,
		"billState":             req.BillState,
		"billZip":               req.BillZip,
		"billCity":              req.BillCity,
		"billAddress":           req.BillAddress,
		"weight":                req.Weight,
		"dimensions":            req.Dimensions,
		"shipCountry":           req.ShipCountry,
		"shipState":             req.ShipState,
		"shipZip":               req.ShipZip,
		"shipCity":              req.ShipCity,
		"shipAddress":           req.ShipAddress,
		"addFraudScore":         req.AddFraudScore,
		"maxPayRetries":         req.MaxPayRetries,
		"reject3dsU":            req.Reject3dsU,
		"payMethod":             req.PayMethod,
		"trType":                req.TrType,
		"extInstallmentoffset":  req.ExtInstallmentOffset,
		"extInstallmentperiod":  req.ExtInstallmentPeriod,
		"extRecurringfrequency": req.ExtRecurringFrequency,
		"extRecurringenddate":   req.ExtRecurringEndDate,
		"blockScore":            req.BlockScore,
		"cssUrl":                req.CssURL,
		"confirmUrl":            req.ConfirmURL,
		"cancelUrl":             req.CancelURL,
		"var1":                  req.Var1,
		"var2":                  req.Var2,
		"var3":                  req.Var3,
		"var4":                  req.Var4,
		"var5":                  req.Var5,
		"var6":                  req.Var6,
		"var7":                  req.Var7,
		"var8":                  req.Var8,
		"var9":                  req.Var9,
		"digest":                dig,
	}
	// Omit empty optional keys for cleaner POST — actually Cardlink may require all keys or only set ones.
	// Docs say form submits hidden inputs; empty optional fields can be omitted from POST.
	// Digest must include ALL positions with "" for omitted — we already computed digest with full preimage.
	// Return only non-empty for optional fields to match typical HTML forms, BUT then browser POST might differ...
	// The doc example includes only some optional fields. The digest is over values merchant sends.
	// So we should return only fields that are non-empty OR required, and for digest we used full 44 positions.
	// Important: if we omit keys from map, client might not send them — that's correct (empty string in digest).
	return pruneEmptyOptional(fields), nil
}

func hasRecurring(req *PaymentRequest) bool {
	return req.ExtRecurringFrequency != "" || req.ExtRecurringEndDate != ""
}

func mutualExclusiveInstallmentRecurring(req *PaymentRequest) error {
	hasInst := req.ExtInstallmentOffset != "" || req.ExtInstallmentPeriod != ""
	hasRec := hasRecurring(req)
	if hasInst && hasRec {
		return fmt.Errorf("redirect: installment and recurring parameters cannot be used together")
	}
	return nil
}

func redirectPreimage(mid string, req *PaymentRequest) string {
	amount := cardlink.FormatOrderAmount(req.OrderAmount)
	return req.Version +
		mid +
		req.Lang +
		req.DeviceCategory +
		req.OrderID +
		req.OrderDesc +
		amount +
		req.Currency +
		req.PayerEmail +
		req.PayerPhone +
		req.BillCountry +
		req.BillState +
		req.BillZip +
		req.BillCity +
		req.BillAddress +
		req.Weight +
		req.Dimensions +
		req.ShipCountry +
		req.ShipState +
		req.ShipZip +
		req.ShipCity +
		req.ShipAddress +
		req.AddFraudScore +
		req.MaxPayRetries +
		req.Reject3dsU +
		req.PayMethod +
		req.TrType +
		req.ExtInstallmentOffset +
		req.ExtInstallmentPeriod +
		req.ExtRecurringFrequency +
		req.ExtRecurringEndDate +
		req.BlockScore +
		req.CssURL +
		req.ConfirmURL +
		req.CancelURL +
		req.Var1 + req.Var2 + req.Var3 + req.Var4 + req.Var5 + req.Var6 + req.Var7 + req.Var8 + req.Var9
}

// pruneEmptyOptional removes empty optional fields from the map; required keys are always kept.
func pruneEmptyOptional(m map[string]string) map[string]string {
	required := map[string]struct{}{
		"version": {}, "mid": {}, "orderid": {}, "orderDesc": {}, "orderAmount": {}, "currency": {},
		"payerEmail": {}, "confirmUrl": {}, "cancelUrl": {}, "digest": {},
	}
	out := make(map[string]string)
	for k, v := range m {
		if _, ok := required[k]; ok || v != "" {
			out[k] = v
		}
	}
	return out
}
