package vposxml

import "errors"

var ErrMissingMID = errors.New("vposxml: MID is required (set Params.MID or Config.MID)")

var ErrMissingSecret = errors.New("vposxml: SharedSecret is required on Config for signed XML requests")

var ErrMissingDigest = errors.New("vposxml: missing or empty Digest in webhook payload")

var ErrDigestMismatch = errors.New("vposxml: webhook digest mismatch")

var errMissingOperation = errors.New("vposxml: Operation is required on RecurringOperationParams")

var ErrMissingTxType = errors.New("vposxml: TxType is required on PaymentLinkParams")
