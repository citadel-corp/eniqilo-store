package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/eniqilo-store/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByPhoneNumberAndUserType(ctx context.Context, phoneNumber string, userType UserType) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// Create implements Repository.
func (d *dbRepository) Create(ctx context.Context, user *User) error {
	createUserQuery := `
		INSERT INTO users (
			id, phone_number, name, user_type, hashed_password
		) VALUES (
			$1, $2, $3, $4, $5
		);
	`
	_, err := d.db.DB().ExecContext(ctx, createUserQuery, user.ID, user.PhoneNumber, user.Name, user.UserType, user.HashedPassword)
	if err != nil {
		return err
	}
	return nil
}

// GetByUsernameAndHashedPassword implements Repository.
func (d *dbRepository) GetByPhoneNumberAndUserType(ctx context.Context, phoneNumber string, userType UserType) (*User, error) {
	getUserQuery := `
		SELECT id, phone_number, name, user_type, hashed_password
		FROM users
		WHERE phone_number = $1 AND user_type = $2;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, phoneNumber, userType)
	u := &User{}
	err := row.Scan(&u.ID, &u.PhoneNumber, &u.Name, &u.UserType, &u.HashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (d *dbRepository) GetByID(ctx context.Context, id string) (*User, error) {
	getUserQuery := `
		SELECT id, phone_number, name, user_type, hashed_password
		FROM users
		WHERE id = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, id)
	u := &User{}
	err := row.Scan(&u.ID, &u.PhoneNumber, &u.Name, &u.UserType, &u.HashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}
