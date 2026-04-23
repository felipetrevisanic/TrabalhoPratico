package interfaces

import "graphql/src/domain/entities"

type ProductRepository interface {
	GetByID(id int) (*entities.Product, error)
	GetAll() ([]entities.Product, error)
	Create(input entities.CreateInput) (*entities.Product, error)
	Update(id int, input entities.UpdateInput) (*entities.Product, error)
	Delete(id int) (bool, error)
	DeleteAll() (bool, error)
}
