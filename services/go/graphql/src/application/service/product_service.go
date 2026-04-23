package service

import (
	appinterfaces "graphql/src/application/interfaces"
	"graphql/src/domain/entities"
	domaininterfaces "graphql/src/domain/interfaces"
)

type ProductService struct {
	repository domaininterfaces.ProductRepository
}

func NewProductService(repository domaininterfaces.ProductRepository) *ProductService {
	return &ProductService{repository: repository}
}

func (s *ProductService) GetByID(id int) (*entities.Product, error) {
	return s.repository.GetByID(id)
}

func (s *ProductService) GetAll() ([]entities.Product, error) {
	return s.repository.GetAll()
}

func (s *ProductService) Create(input entities.CreateInput) (*entities.Product, error) {
	return s.repository.Create(input)
}

func (s *ProductService) Update(id int, input entities.UpdateInput) (*entities.Product, error) {
	return s.repository.Update(id, input)
}

func (s *ProductService) Delete(id int) (bool, error) {
	return s.repository.Delete(id)
}

func (s *ProductService) DeleteAll() (bool, error) {
	return s.repository.DeleteAll()
}

var _ appinterfaces.ProductService = (*ProductService)(nil)
