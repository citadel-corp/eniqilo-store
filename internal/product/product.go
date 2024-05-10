package product

import "time"

type ProductCategory string

var (
	CategoryClothing    ProductCategory = "Clothing"
	CategoryFootwear    ProductCategory = "Footwear"
	CategoryAccessories ProductCategory = "Accessories"
	CategoryBeverages   ProductCategory = "Beverages"
)

var ProductCategories = []interface{}{CategoryAccessories, CategoryFootwear, CategoryAccessories, CategoryBeverages}

type Product struct {
	ID          string
	UserID      string
	Name        string
	SKU         string
	Category    ProductCategory
	ImageURL    string
	Notes       string
	Price       int64
	Stock       int
	Location    string
	IsAvailable bool
	CreatedAt   time.Time
}
