package product

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

type Repository interface {
	GetByID(id int) (*Product, error)
	GetAll() ([]Product, error)
	Create(input CreateInput) (*Product, error)
	Update(id int, input UpdateInput) (*Product, error)
	Delete(id int) (bool, error)
	DeleteAll() (bool, error)
}
