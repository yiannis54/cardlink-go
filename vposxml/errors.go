package vposxml

import "errors"

var ErrMissingMID = errors.New("vposxml: MID is required (set Params.MID or Config.MID)")

var ErrMissingSecret = errors.New("vposxml: SharedSecret is required on Config for signed XML requests")

var errMissingOperation = errors.New("vposxml: Operation is required on RecurringOperationParams")
