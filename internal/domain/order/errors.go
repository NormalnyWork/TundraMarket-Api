package order

import "errors"

var (
	ErrEmptyCart      = errors.New("EMPTY_CART")
	ErrInvalidId      = errors.New("INVALID_ID")
	ErrDistanceTooFar = errors.New("DISTANCE_TOO_FAR")
)
