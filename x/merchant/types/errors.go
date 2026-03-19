package types

import "errors"

var (
	ErrMaxMerchantsReached = errors.New("max merchants reached")
	ErrMerchantNotFound    = errors.New("merchant not found")
	ErrMerchantAlreadyExists = errors.New("merchant already exists")
	ErrNotMerchantOwner    = errors.New("not merchant owner")
	ErrInvalidApiKey       = errors.New("invalid api key")
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrInvalidPaymentStatus = errors.New("invalid payment status")
)
