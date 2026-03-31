package product

import "time"

type Product struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Price         float64    `json:"price"`
	StockQuantity int        `json:"stockQuantity"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
}

type CreateRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stockQuantity"`
}

type UpdateRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stockQuantity"`
}

type Repository interface {
	GetByID(id int) (*Product, error)
	GetAll() ([]Product, error)
	Create(request CreateRequest) (*Product, error)
	Update(id int, request UpdateRequest) (*Product, error)
	Delete(id int) (bool, error)
	DeleteAll() (bool, error)
}
