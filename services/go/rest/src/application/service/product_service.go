package service

import (
	appinterfaces "rest/src/application/interfaces"
	"rest/src/domain/entities"
	domaininterfaces "rest/src/domain/interfaces"
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

func (s *ProductService) Create(request entities.CreateRequest) (*entities.Product, error) {
	return s.repository.Create(request)
}

func (s *ProductService) Update(id int, request entities.UpdateRequest) (*entities.Product, error) {
	return s.repository.Update(id, request)
}

func (s *ProductService) Delete(id int) (bool, error) {
	return s.repository.Delete(id)
}

func (s *ProductService) DeleteAll() (bool, error) {
	return s.repository.DeleteAll()
}

var _ appinterfaces.ProductService = (*ProductService)(nil)
