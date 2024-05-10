package product

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var imgUrlValidationRule = validation.NewStringRule(func(s string) bool {
	match, _ := regexp.MatchString(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/{1}[A-z0-9_\-\:x\=\(\)]+)*(\.(jpg|jpeg|png))?$`, s)
	return match
}, "image url is not valid")

type CreateProductPayload struct {
	Name        string          `json:"name"`
	SKU         string          `json:"sku"`
	Category    ProductCategory `json:"category"`
	ImageURL    string          `json:"imageURL"`
	Notes       string          `json:"notes"`
	Price       int64           `json:"price"`
	Stock       int             `json:"stock"`
	Location    string          `json:"location"`
	IsAvailable bool            `json:"isAvailable"`
}

func (p CreateProductPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required, validation.Length(1, 30)),
		validation.Field(&p.SKU, validation.Required, validation.Length(1, 30)),
		validation.Field(&p.Category, validation.Required, validation.In(ProductCategories...)),
		validation.Field(&p.ImageURL, validation.Required, imgUrlValidationRule),
		validation.Field(&p.Notes, validation.Required, validation.Length(1, 200)),
		validation.Field(&p.Price, validation.Required, validation.Min(1)),
		validation.Field(&p.Stock, validation.Required, validation.Min(1)),
		validation.Field(&p.Location, validation.Required, validation.Length(1, 200)),
		validation.Field(&p.IsAvailable, validation.Required),
	)
}
