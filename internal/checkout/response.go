package checkout

import "time"

type CheckoutHistoryResponse struct {
	TransactionID  string                  `json:"transactionId"`
	CustomerID     string                  `json:"customerId"`
	ProductDetails []ProductDetailResponse `json:"productDetails"`
	Paid           int                     `json:"paid"`
	Change         int                     `json:"change"`
	CreatedAt      time.Time               `json:"createdAt"`
}

type ProductDetailResponse struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
