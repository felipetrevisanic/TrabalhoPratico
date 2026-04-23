package interfaces

import "rest/src/domain/entities"

type ProductService interface {
	GetByID(id int) (*entities.Product, error)
	GetAll() ([]entities.Product, error)
	Create(request entities.CreateRequest) (*entities.Product, error)
	Update(id int, request entities.UpdateRequest) (*entities.Product, error)
	Delete(id int) (bool, error)
	DeleteAll() (bool, error)
}
