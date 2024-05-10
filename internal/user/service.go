package user

import (
	"context"
	"fmt"
	"time"

	"github.com/citadel-corp/eniqilo-store/internal/common/id"
	"github.com/citadel-corp/eniqilo-store/internal/common/jwt"
	"github.com/citadel-corp/eniqilo-store/internal/common/password"
)

type Service interface {
	CreateStaff(ctx context.Context, req CreateStaffPayload) (*StaffResponse, error)
	CreateCustomer(ctx context.Context, req CreateCustomerPayload) (*CustomerResponse, error)
	StaffLogin(ctx context.Context, req LoginPayload) (*StaffResponse, error)
	ListCustomers(ctx context.Context, req ListCustomerPayload) ([]*CustomerResponse, error)
}

type userService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}

func (s *userService) CreateStaff(ctx context.Context, req CreateStaffPayload) (*StaffResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.repository.GetByPhoneNumberAndUserType(ctx, req.PhoneNumber, Staff)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, ErrPhoneNumberAlreadyExists
	}
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	user = &User{
		ID:             id.GenerateStringID(16),
		UserType:       Staff,
		PhoneNumber:    req.PhoneNumber,
		Name:           req.Name,
		HashedPassword: hashedPassword,
	}
	err = s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Hour*2, fmt.Sprint(user.ID))
	if err != nil {
		return nil, err
	}
	return &StaffResponse{
		UserID:      user.ID,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		AccessToken: accessToken,
	}, nil
}

func (s *userService) CreateCustomer(ctx context.Context, req CreateCustomerPayload) (*CustomerResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.repository.GetByPhoneNumberAndUserType(ctx, req.PhoneNumber, Customer)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, ErrPhoneNumberAlreadyExists
	}
	user = &User{
		ID:          id.GenerateStringID(16),
		UserType:    Staff,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
	}
	err = s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &CustomerResponse{
		UserID:      user.ID,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
	}, nil
}

func (s *userService) StaffLogin(ctx context.Context, req LoginPayload) (*StaffResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.repository.GetByPhoneNumberAndUserType(ctx, req.PhoneNumber, Staff)
	if err != nil {
		return nil, err
	}
	match, err := password.Matches(req.Password, user.HashedPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrWrongPassword
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Hour*2, fmt.Sprint(user.ID))
	if err != nil {
		return nil, err
	}
	return &StaffResponse{
		UserID:      user.ID,
		PhoneNumber: user.PhoneNumber,
		Name:        user.Name,
		AccessToken: accessToken,
	}, nil
}

// ListCustomer implements Service.
func (s *userService) ListCustomers(ctx context.Context, req ListCustomerPayload) ([]*CustomerResponse, error) {
	customers, err := s.repository.ListCustomers(ctx, req)
	if err != nil {
		return nil, err
	}
	res := make([]*CustomerResponse, len(customers))
	for i, customer := range customers {
		res[i] = &CustomerResponse{
			UserID:      customer.ID,
			PhoneNumber: customer.PhoneNumber,
			Name:        customer.Name,
		}
	}
	return res, nil
}
