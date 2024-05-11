package product

import "errors"

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrProductNotFound  = errors.New("product not found")
)
