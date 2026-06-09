package order

import "errors"

var (
	ErrEmptyCart           = errors.New("EMPTY_CART")
	ErrInvalidId           = errors.New("INVALID_ID")
	ErrDistanceTooFar      = errors.New("DISTANCE_TOO_FAR")
	ErrUnknownStatus       = errors.New("UNKNOWN_STATUS")
	ErrUnknownCategory     = errors.New("UNKNOWN_CATEGORY")
	ErrIllegalStatusChange = errors.New("ILLEGAL_STATUS_CHANGE")
	ErrForbidden           = errors.New("FORBIDDEN")
	ErrInvalidPhone        = errors.New("INVALID_PHONE")
	ErrNomadNotFound       = errors.New("NOMAD_NOT_FOUND")
)
