package checkout

import (
	"bytes"
	"context"
	"fmt"

	"github.com/citadel-corp/eniqilo-store/internal/common/db"
)

type Repository interface {
	ListCheckoutHistories(ctx context.Context, req ListCheckoutHistoriesPayload) ([]*CheckoutHistory, error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// ListCustomers implements Repository.
func (d *dbRepository) ListCheckoutHistories(ctx context.Context, req ListCheckoutHistoriesPayload) ([]*CheckoutHistory, error) {
	var query bytes.Buffer
	_, _ = query.WriteString("SELECT * FROM checkout_histories ")
	params := make([]interface{}, 0)
	if req.CustomerID != "" {
		_, _ = query.WriteString("WHERE user_id = $1 ")
		params = append(params, req.CustomerID)
	}
	switch req.CreatedAtSearchType {
	case Ascending:
		_, _ = query.WriteString(" ORDER BY created_at ASC ")
	case Descending:
		_, _ = query.WriteString(" ORDER BY created_at DESC ")
	}
	_, _ = query.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d;", req.Limit, req.Offset))
	rows, err := d.db.DB().QueryContext(ctx, query.String(), params...)
	if err != nil {
		return nil, err
	}
	res := make([]*CheckoutHistory, 0)
	for rows.Next() {
		ch := &CheckoutHistory{}
		err := rows.Scan(&ch.ID, &ch.UserID, &ch.ProductDetails, &ch.Paid, &ch.Change, &ch.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, ch)
	}
	return res, nil
}
