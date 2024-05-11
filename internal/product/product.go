package product

import (
	"errors"
	"time"
)

type ProductCategory string

var (
	CategoryClothing    ProductCategory = "Clothing"
	CategoryFootwear    ProductCategory = "Footwear"
	CategoryAccessories ProductCategory = "Accessories"
	CategoryBeverages   ProductCategory = "Beverages"
)

var ProductCategories = []interface{}{CategoryClothing, CategoryFootwear, CategoryAccessories, CategoryBeverages}

func ParseProductCategory(str string) (ProductCategory, error) {
	for i := range ProductCategories {
		if ProductCategory(str) == ProductCategories[i] {
			return ProductCategory(str), nil
		}
	}
	return ProductCategory(""), errors.New("not a product category")
}

type Product struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	SKU         string          `json:"sku"`
	Category    ProductCategory `json:"category"`
	ImageURL    string          `json:"imageUrl"`
	Notes       string          `json:"notes"`
	Price       int64           `json:"price"`
	Stock       int             `json:"stock"`
	Location    string          `json:"location"`
	IsAvailable bool            `json:"isAvailable"`
	CreatedAt   time.Time       `json:"createdAt"`
}
