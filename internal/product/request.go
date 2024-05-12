package product

import (
	"errors"
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
	Stock       *int            `json:"stock"`
	Location    string          `json:"location"`
	IsAvailable *bool           `json:"isAvailable"`
}

func (p CreateProductPayload) Validate() error {
	if p.IsAvailable == nil {
		return errors.New("isAvailable: required field")
	}
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required, validation.Length(1, 30)),
		validation.Field(&p.SKU, validation.Required, validation.Length(1, 30)),
		validation.Field(&p.Category, validation.Required, validation.In(ProductCategories...)),
		validation.Field(&p.ImageURL, validation.Required, imgUrlValidationRule),
		validation.Field(&p.Notes, validation.Required, validation.Length(1, 200)),
		validation.Field(&p.Price, validation.Required, validation.Min(1)),
		validation.Field(&p.Stock, validation.NotNil, validation.Min(0), validation.Max(100000)),
		validation.Field(&p.Location, validation.Required, validation.Length(1, 200)),
	)
}

type EditProductPayload struct {
	ID          string          `json:"-"`
	Name        string          `json:"name"`
	SKU         string          `json:"sku"`
	Category    ProductCategory `json:"category"`
	ImageURL    string          `json:"imageURL"`
	Notes       string          `json:"notes"`
	Price       int64           `json:"price"`
	Stock       int             `json:"stock"`
	Location    string          `json:"location"`
	IsAvailable *bool           `json:"isAvailable"`
}

func (p EditProductPayload) Validate() error {
	if p.IsAvailable == nil {
		return errors.New("isAvailable: required field")
	}
	return validation.ValidateStruct(&p,
		validation.Field(&p.ID, validation.Required),
		validation.Field(&p.Name, validation.Required, validation.Length(1, 30)),
		validation.Field(&p.SKU, validation.Required, validation.Length(1, 30)),
		validation.Field(&p.Category, validation.Required, validation.In(ProductCategories...)),
		validation.Field(&p.ImageURL, validation.Required, imgUrlValidationRule),
		validation.Field(&p.Notes, validation.Required, validation.Length(1, 200)),
		validation.Field(&p.Price, validation.Required, validation.Min(1)),
		validation.Field(&p.Stock, validation.Required, validation.Min(1), validation.Max(100000)),
		validation.Field(&p.Location, validation.Required, validation.Length(1, 200)),
	)
}

type DeleteProductPayload struct {
	ID string
}

func (p DeleteProductPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ID, validation.Required),
	)
}

type ListProductPayload struct {
	ID          string `schema:"id" binding:"omitempty"`
	Limit       int    `schema:"limit" binding:"omitempty"`
	Offset      int    `schema:"offset" binding:"omitempty"`
	Name        string `schema:"name" binding:"omitempty"`
	IsAvailable string `schema:"isAvailable" binding:"omitempty"`
	Category    string `schema:"category" binding:"omitempty"`
	SKU         string `schema:"sku" binding:"omitempty"`
	Price       string `schema:"price" binding:"omitempty"`
	InStock     string `schema:"inStock" binding:"omitempty"`
	CreatedAt   string `schema:"createdAt" binding:"omitempty"`
}
