package server

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	appinterfaces "grpc/src/application/interfaces"
	"grpc/src/domain/entities"
	productv1 "grpc/src/interfaces/grpc/gen/productv1"
)

type ProductServer struct {
	productv1.UnimplementedProductServiceServer
	service appinterfaces.ProductService
}

func NewProductServer(service appinterfaces.ProductService) *ProductServer {
	return &ProductServer{service: service}
}

func (s *ProductServer) GetProductById(ctx context.Context, request *productv1.GetProductByIdRequest) (*productv1.ProductResponse, error) {
	item, err := s.service.GetByID(int(request.Id))
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return mapProduct(*item), nil
}

func (s *ProductServer) GetAllProducts(ctx context.Context, request *productv1.GetAllProductsRequest) (*productv1.GetAllProductsResponse, error) {
	items, err := s.service.GetAll()
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
	item, err := s.service.Create(entities.CreateInput{
		Name:          request.Name,
		Description:   request.Description,
		Category:      request.Category,
		Images:        request.Images,
		Price:         request.Price,
		StockQuantity: int(request.StockQuantity),
	})
	if err != nil {
		return nil, err
	}

	return mapProduct(*item), nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, request *productv1.UpdateProductRequest) (*productv1.ProductResponse, error) {
	item, err := s.service.Update(int(request.Id), entities.UpdateInput{
		Name:          request.Name,
		Description:   request.Description,
		Category:      request.Category,
		Images:        request.Images,
		Price:         request.Price,
		StockQuantity: int(request.StockQuantity),
	})
	if err != nil {
		return nil, err
	}

	return mapProduct(*item), nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, request *productv1.DeleteProductRequest) (*productv1.DeleteProductResponse, error) {
	deleted, err := s.service.Delete(int(request.Id))
	if err != nil {
		return nil, err
	}

	return &productv1.DeleteProductResponse{Deleted: deleted}, nil
}

func (s *ProductServer) DeleteAllProducts(ctx context.Context, request *productv1.DeleteAllProductsRequest) (*productv1.DeleteAllProductsResponse, error) {
	deleted, err := s.service.DeleteAll()
	if err != nil {
		return nil, err
	}

	return &productv1.DeleteAllProductsResponse{Deleted: deleted}, nil
}

func mapProduct(item entities.Product) *productv1.ProductResponse {
	response := &productv1.ProductResponse{
		Id:            int32(item.ID),
		Name:          item.Name,
		Description:   item.Description,
		Category:      item.Category,
		Images:        item.Images,
		Price:         item.Price,
		StockQuantity: int32(item.StockQuantity),
		CreatedAt:     item.CreatedAt.UTC().Format(time.RFC3339Nano),
	}

	if item.UpdatedAt != nil {
		response.UpdatedAt = item.UpdatedAt.UTC().Format(time.RFC3339Nano)
	}

	return response
}
