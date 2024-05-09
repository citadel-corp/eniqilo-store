package user

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var phoneNumberValidationRule = validation.NewStringRule(func(s string) bool {
	return strings.HasPrefix(s, "+")
}, "phone number must start with international calling code")

type CreateUserPayload struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}

func (p CreateUserPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.PhoneNumber, validation.Required, phoneNumberValidationRule, validation.Length(10, 16)),
		validation.Field(&p.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 15)),
	)
}

type LoginPayload struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

func (p LoginPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.PhoneNumber, validation.Required, phoneNumberValidationRule, validation.Length(10, 16)),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 15)),
	)
}
