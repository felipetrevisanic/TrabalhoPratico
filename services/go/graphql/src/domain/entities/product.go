package entities

import "time"

type Product struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Category      string     `json:"category"`
	Images        []string   `json:"images"`
	Price         float64    `json:"price"`
	StockQuantity int        `json:"stockQuantity"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
}

type CreateInput struct {
	Name          string
	Description   string
	Category      string
	Images        []string
	Price         float64
	StockQuantity int
}

type UpdateInput struct {
	Name          string
	Description   string
	Category      string
	Images        []string
	Price         float64
	StockQuantity int
}
