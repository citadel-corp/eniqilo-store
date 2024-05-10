package product

import (
	"context"

	"github.com/citadel-corp/eniqilo-store/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	GetByMultipleID(ctx context.Context, ids []string) ([]*Product, error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

func (d *dbRepository) Create(ctx context.Context, product *Product) (*Product, error) {
	createUserQuery := `
		INSERT INTO products (
			id, name, sku, category, image_url, notes, price, stock, location, is_available
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING created_at;
	`
	row := d.db.DB().QueryRowContext(ctx, createUserQuery,
		product.ID, product.Name, product.SKU, product.Category, product.ImageURL, product.Notes, product.Price, product.Stock,
		product.Location, product.IsAvailable)
	err := row.Scan(&product.CreatedAt)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (d *dbRepository) GetByMultipleID(ctx context.Context, ids []string) ([]*Product, error) {
	q := `
		SELECT *
		FROM products
		WHERE id IN ?;
	`
	rows, err := d.db.DB().QueryContext(ctx, q, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]*Product, 0)
	for rows.Next() {
		p := &Product{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.SKU, &p.Category, &p.ImageURL, &p.Notes, &p.Price, &p.Stock, &p.Location, &p.IsAvailable, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}
