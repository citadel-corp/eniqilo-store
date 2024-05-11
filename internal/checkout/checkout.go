package checkout

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type CheckoutHistory struct {
	ID             string
	UserID         string
	ProductDetails ProductDetails
	Paid           int
	Change         int
	CreatedAt      time.Time
}

type ProductDetail struct {
	ProductID string
	Quantity  int
}

type ProductDetails []ProductDetail

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a ProductDetails) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *ProductDetails) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
