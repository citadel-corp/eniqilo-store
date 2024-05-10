package checkout

import (
	"context"
)

type Service interface {
	ListCheckoutHistories(ctx context.Context, req ListCheckoutHistoriesPayload) ([]*CheckoutHistoryResponse, error)
}

type checkoutService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &checkoutService{repository: repository}
}

func (s *checkoutService) ListCheckoutHistories(ctx context.Context, req ListCheckoutHistoriesPayload) ([]*CheckoutHistoryResponse, error) {
	req.CreatedAtSearchType = Ignore
	switch req.CreatedAt {
	case "asc":
		req.CreatedAtSearchType = Ascending
	case "desc":
		req.CreatedAtSearchType = Descending
	}
	checkoutHistories, err := s.repository.ListCheckoutHistories(ctx, req)
	if err != nil {
		return nil, err
	}
	res := make([]*CheckoutHistoryResponse, len(checkoutHistories))
	for i, checkoutHistory := range checkoutHistories {
		productDetails := make([]ProductDetailResponse, len(checkoutHistory.ProductDetails))
		for j, productDetail := range checkoutHistory.ProductDetails {
			productDetails[j] = ProductDetailResponse{
				ProductID: productDetail.ProductID,
				Quantity:  productDetail.Quantity,
			}
		}
		res[i] = &CheckoutHistoryResponse{
			TransactionID:  checkoutHistory.ID,
			CustomerID:     checkoutHistory.UserID,
			ProductDetails: productDetails,
			Paid:           checkoutHistory.Paid,
			Change:         checkoutHistory.Change,
			CreatedAt:      checkoutHistory.CreatedAt,
		}
	}
	return res, nil
}
