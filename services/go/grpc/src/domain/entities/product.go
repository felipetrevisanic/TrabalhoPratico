package entities

import "time"

type Product struct {
	ID            int
	Name          string
	Description   string
	Price         float64
	StockQuantity int
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type CreateInput struct {
	Name          string
	Description   string
	Price         float64
	StockQuantity int
}

type UpdateInput struct {
	Name          string
	Description   string
	Price         float64
	StockQuantity int
}
