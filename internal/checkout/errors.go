package checkout

import "errors"

var (
	ErrCustomerNotFound      = errors.New("customer id is not found")
	ErrProductNotFound       = errors.New("one or more products is not available")
	ErrProductUnavailable    = errors.New("product is unavailable")
	ErrProductStockNotEnough = errors.New("product stock is not enough")
	ErrNotEnoughMoney        = errors.New("not enough money paid")
	ErrWrongChange           = errors.New("wrong change")
	ErrValidationFailed      = errors.New("validation failed")
)
