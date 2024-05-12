package product

import (
	"context"
	"fmt"
	"strconv"

	"github.com/citadel-corp/eniqilo-store/internal/common/db"
)

type Repository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	GetByMultipleID(ctx context.Context, ids []string) ([]*Product, error)
	Put(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, req ListProductPayload) ([]Product, error)
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
	if len(ids) == 0 {
		return make([]*Product, 0), nil
	}
	q := `
		SELECT *
		FROM products
		WHERE id IN(
	`

	for i, v := range ids {
		if i > 0 {
			q += ","
		}
		q += fmt.Sprintf("'%s'", v)

	}

	q += ");"
	rows, err := d.db.DB().QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]*Product, 0)
	for rows.Next() {
		p := &Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Category, &p.ImageURL, &p.Notes, &p.Price, &p.Stock, &p.Location, &p.IsAvailable, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (d *dbRepository) Put(ctx context.Context, product *Product) error {
	q := `
        UPDATE products
        SET name = $1, sku = $2, category = $3, image_url = $4, notes = $5, price = $6, stock = $7, location = $8, is_available = $9
        WHERE id = $10;
    `
	row, err := d.db.DB().ExecContext(ctx, q, product.Name, product.SKU, product.Category, product.ImageURL, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable, product.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (d *dbRepository) Delete(ctx context.Context, id string) error {
	q := `
        DELETE FROM products
        WHERE id = $1;
    `
	row, err := d.db.DB().ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (d *dbRepository) List(ctx context.Context, req ListProductPayload) ([]Product, error) {
	q := `
		SELECT id, name, sku, category, image_url, stock, notes, price, location, is_available, created_at
		FROM products
	`
	paramNo := 1
	params := make([]interface{}, 0)
	if req.ID != "" {
		q += fmt.Sprintf("WHERE id = $%d ", paramNo)
		paramNo += 1
		params = append(params, req.ID)
	}
	if req.Name != "" {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("LOWER(name) LIKE $%d ", paramNo)
		paramNo += 1
		params = append(params, "%"+req.Name+"%")
	}
	if v, err := strconv.ParseBool(req.IsAvailable); err == nil {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("is_available = $%d ", paramNo)
		paramNo += 1
		params = append(params, v)
	}
	if v, err := ParseProductCategory(req.Category); err == nil {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("category = $%d ", paramNo)
		paramNo += 1
		params = append(params, v)
	}
	if req.SKU != "" {
		q += whereOrAnd(paramNo)
		q += fmt.Sprintf("sku = $%d ", paramNo)
		paramNo += 1
		params = append(params, req.SKU)
	}

	if v, err := strconv.ParseBool(req.InStock); err == nil {
		q += whereOrAnd(paramNo)
		if v {
			q += "stock > 0 "
		} else {
			q += "stock = 0 "
		}
	}

	var orderedByPrice bool
	if req.Price != "" {
		if req.Price == "asc" || req.Price == "desc" {
			orderedByPrice = true
			orderBy := "asc"
			if req.Price == "desc" {
				orderBy = "desc"
			}

			q += `ORDER BY price ` + orderBy
		}
	}

	orderByCreatedAt := "desc"
	if req.CreatedAt == "asc" || req.CreatedAt == "desc" {
		if req.CreatedAt == "asc" {
			orderByCreatedAt = "asc"
		}
	}

	if orderedByPrice {
		q += `, created_at ` + orderByCreatedAt
	} else {
		q += `ORDER BY created_at ` + orderByCreatedAt
	}

	q += fmt.Sprintf(" OFFSET $%d LIMIT $%d", paramNo, paramNo+1)
	params = append(params, req.Offset)
	params = append(params, req.Limit)

	rows, err := d.db.DB().QueryContext(ctx, q, params...)
	if err != nil {
		return nil, err
	}
	res := make([]Product, 0)
	for rows.Next() {
		product := Product{}
		err = rows.Scan(&product.ID, &product.Name, &product.SKU, &product.Category, &product.ImageURL, &product.Stock,
			&product.Notes, &product.Price, &product.Location, &product.IsAvailable, &product.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, product)
	}
	return res, nil
}

func whereOrAnd(paramNo int) string {
	if paramNo == 1 {
		return "WHERE "
	}
	return "AND "
}
