package checkout

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CheckoutRequest struct {
	CustomerID     string                 `json:"customerId"`
	ProductDetails []ProductDetailRequest `json:"productDetails"`
	Paid           int                    `json:"paid"`
	Change         *int                   `json:"change"`
}

func (p CheckoutRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.CustomerID, validation.Required),
		validation.Field(&p.Paid, validation.Required, validation.Min(1)),
		validation.Field(&p.Change, validation.NotNil),
	)
}

type ProductDetailRequest struct {
	ProductID     string `json:"productId"`
	Quantity      int    `json:"quantity"`
	OriginalStock int    `json:"originalStock"`
}

func (p ProductDetailRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ProductID, validation.Required),
		validation.Field(&p.Quantity, validation.Required, validation.Min(1)),
	)
}

type ListCheckoutHistoriesPayload struct {
	CustomerID string `schema:"customerId" binding:"omitempty"`
	Limit      int    `schema:"limit" binding:"omitempty"`
	Offset     int    `schema:"offset" binding:"omitempty"`
	CreatedAt  string `schema:"createdAt" binding:"omitempty"`

	CreatedAtSearchType CreatedAtSearchType
}

type CreatedAtSearchType int

const (
	Ascending CreatedAtSearchType = iota
	Descending
	Ignore
)
