package product

import (
	"context"

	"github.com/citadel-corp/eniqilo-store/internal/common/id"
)

type Service interface {
	Create(ctx context.Context, req CreateProductPayload) (*ProductResponse, error)
	Edit(ctx context.Context, req EditProductPayload) error
	Delete(ctx context.Context, req DeleteProductPayload) error
	List(ctx context.Context, req ListProductPayload) ([]Product, error)
	ListForCustomers(ctx context.Context, req ListProductPayload) ([]Product, error)
}

type productService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &productService{repository: repository}
}

func (s *productService) Create(ctx context.Context, req CreateProductPayload) (*ProductResponse, error) {
	product := &Product{
		ID:          id.GenerateStringID(16),
		Name:        req.Name,
		SKU:         req.SKU,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		Notes:       req.Notes,
		Price:       req.Price,
		Stock:       req.Stock,
		Location:    req.Location,
		IsAvailable: *req.IsAvailable,
	}

	product, err := s.repository.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return &ProductResponse{
		ID:        product.ID,
		CreatedAt: product.CreatedAt,
	}, nil
}

func (s *productService) Edit(ctx context.Context, req EditProductPayload) error {
	product := &Product{
		ID:          req.ID,
		Name:        req.Name,
		SKU:         req.SKU,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		Notes:       req.Notes,
		Price:       req.Price,
		Stock:       req.Stock,
		Location:    req.Location,
		IsAvailable: *req.IsAvailable,
	}
	err := s.repository.Put(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) Delete(ctx context.Context, req DeleteProductPayload) error {
	err := s.repository.Delete(ctx, req.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) List(ctx context.Context, req ListProductPayload) ([]Product, error) {
	if req.Limit == 0 {
		req.Limit = 5
	}

	products, err := s.repository.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *productService) ListForCustomers(ctx context.Context, req ListProductPayload) ([]Product, error) {
	if req.Limit == 0 {
		req.Limit = 5
	}

	req.IsAvailable = "true"
	products, err := s.repository.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return products, nil
}
