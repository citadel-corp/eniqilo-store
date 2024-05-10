package checkout

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrWrongPassword            = errors.New("wrong password")
	ErrPhoneNumberAlreadyExists = errors.New("phone number already exists")
	ErrValidationFailed         = errors.New("validation failed")
)
