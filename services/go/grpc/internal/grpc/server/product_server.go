package server

import (
	"context"
	"fmt"
	"time"

	productv1 "grpc/internal/grpc/gen/productv1"
	"grpc/internal/product"
)

type ProductServer struct {
	productv1.UnimplementedProductServiceServer
	repository product.Repository
}

func NewProductServer(repository product.Repository) *ProductServer {
	return &ProductServer{repository: repository}
}

func (s *ProductServer) GetProductById(ctx context.Context, request *productv1.GetProductByIdRequest) (*productv1.ProductResponse, error) {
	item, err := s.repository.GetByID(int(request.Id))
	if err != nil {
		return nil, err
	}

	if item == nil {
		item = &product.Product{
			ID:            int(request.Id),
			Name:          fmt.Sprintf("Product %d", request.Id),
			Description:   "Product not found in sample list",
			Price:         0,
			StockQuantity: 0,
			CreatedAt:     time.Now().UTC(),
		}
	}

	return mapProduct(*item), nil
}

func (s *ProductServer) GetAllProducts(ctx context.Context, request *productv1.GetAllProductsRequest) (*productv1.GetAllProductsResponse, error) {
	items, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	response := &productv1.GetAllProductsResponse{
		Products: make([]*productv1.ProductResponse, 0, len(items)),
	}

	for _, item := range items {
		response.Products = append(response.Products, mapProduct(item))
	}

	return response, nil
}

func (s *ProductServer) CreateProduct(ctx context.Context, request *productv1.CreateProductRequest) (*productv1.ProductResponse, error) {
	item, err := s.repository.Create(product.CreateInput{
		Name:          request.Name,
		Description:   request.Description,
		Price:         request.Price,
		StockQuantity: int(request.StockQuantity),
	})
	if err != nil {
		return nil, err
	}

	return mapProduct(*item), nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, request *productv1.UpdateProductRequest) (*productv1.ProductResponse, error) {
	item, err := s.repository.Update(int(request.Id), product.UpdateInput{
		Name:          request.Name,
		Description:   request.Description,
		Price:         request.Price,
		StockQuantity: int(request.StockQuantity),
	})
	if err != nil {
		return nil, err
	}

	return mapProduct(*item), nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, request *productv1.DeleteProductRequest) (*productv1.DeleteProductResponse, error) {
	deleted, err := s.repository.Delete(int(request.Id))
	if err != nil {
		return nil, err
	}

	return &productv1.DeleteProductResponse{Deleted: deleted}, nil
}

func (s *ProductServer) DeleteAllProducts(ctx context.Context, request *productv1.DeleteAllProductsRequest) (*productv1.DeleteAllProductsResponse, error) {
	deleted, err := s.repository.DeleteAll()
	if err != nil {
		return nil, err
	}

	return &productv1.DeleteAllProductsResponse{Deleted: deleted}, nil
}

func mapProduct(item product.Product) *productv1.ProductResponse {
	response := &productv1.ProductResponse{
		Id:            int32(item.ID),
		Name:          item.Name,
		Description:   item.Description,
		Price:         item.Price,
		StockQuantity: int32(item.StockQuantity),
		CreatedAt:     item.CreatedAt.UTC().Format(time.RFC3339Nano),
	}

	if item.UpdatedAt != nil {
		response.UpdatedAt = item.UpdatedAt.UTC().Format(time.RFC3339Nano)
	}

	return response
}
