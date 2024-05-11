package checkout

import (
	"context"
	"fmt"

	"github.com/citadel-corp/eniqilo-store/internal/common/id"
	"github.com/citadel-corp/eniqilo-store/internal/product"
	"github.com/citadel-corp/eniqilo-store/internal/user"
)

type Service interface {
	CheckoutProducts(ctx context.Context, req CheckoutRequest) error
	ListCheckoutHistories(ctx context.Context, req ListCheckoutHistoriesPayload) ([]*CheckoutHistoryResponse, error)
}

type checkoutService struct {
	repository        Repository
	userRepository    user.Repository
	productRepository product.Repository
}

func NewService(repository Repository, userRepository user.Repository, productRepository product.Repository) Service {
	return &checkoutService{
		repository:        repository,
		userRepository:    userRepository,
		productRepository: productRepository,
	}
}

// CheckoutProducts implements Service.
func (s *checkoutService) CheckoutProducts(ctx context.Context, req CheckoutRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	productIDs := make([]string, len(req.ProductDetails))
	productMap := make(map[string]ProductDetailRequest, len(req.ProductDetails))
	productDetails := make([]ProductDetail, len(req.ProductDetails))
	for i, productDetail := range req.ProductDetails {
		if err := productDetail.Validate(); err != nil {
			return fmt.Errorf("%w: %w", ErrValidationFailed, err)
		}
		productIDs[i] = productDetail.ProductID
		productMap[productDetail.ProductID] = productDetail
		productDetails[i] = ProductDetail{
			ProductID: productDetail.ProductID,
			Quantity:  productDetail.Quantity,
		}
	}

	user, err := s.userRepository.GetByID(ctx, req.CustomerID)
	if err != nil {
		return err
	}
	products, err := s.productRepository.GetByMultipleID(ctx, productIDs)
	if err != nil {
		return err
	}
	if len(productIDs) != len(products) {
		return ErrProductNotFound
	}

	price := int64(0)
	for _, product := range products {
		if !product.IsAvailable {
			return ErrProductUnavailable
		}
		if product.Stock-productMap[product.ID].Quantity < 0 {
			return ErrProductStockNotEnough
		}
		price += product.Price * int64(productMap[product.ID].Quantity)
	}
	if int64(req.Paid) < price {
		return ErrNotEnoughMoney
	}
	change := req.Paid - int(price)
	if change != *req.Change {
		return ErrWrongChange
	}
	ch := &CheckoutHistory{
		ID:             id.GenerateStringID(16),
		UserID:         user.ID,
		ProductDetails: productDetails,
		Paid:           req.Paid,
		Change:         *req.Change,
	}
	return s.repository.CreateCheckoutHistory(ctx, ch)
}

func (s *checkoutService) ListCheckoutHistories(ctx context.Context, req ListCheckoutHistoriesPayload) ([]*CheckoutHistoryResponse, error) {
	req.CreatedAtSearchType = Descending
	switch req.CreatedAt {
	case "asc":
		req.CreatedAtSearchType = Ascending
	case "desc":
		req.CreatedAtSearchType = Descending
	}
	if req.Limit == 0 {
		req.Limit = 5
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
