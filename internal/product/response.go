package product

import "time"

type ProductResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
