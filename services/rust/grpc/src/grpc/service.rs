use std::sync::Arc;

use chrono::Utc;
use rust_decimal::{
    Decimal,
    prelude::{FromPrimitive, ToPrimitive},
};
use tonic::{Request, Response, Status};

use crate::{
    grpc::productv1::{
        CreateProductRequest, DeleteAllProductsRequest, DeleteAllProductsResponse,
        DeleteProductRequest, DeleteProductResponse, GetAllProductsRequest, GetAllProductsResponse,
        GetProductByIdRequest, ProductResponse, UpdateProductRequest,
        product_service_server::ProductService,
    },
    product::{CreateProductInput, Product, UpdateProductInput},
    repository::product_repository::{PostgresProductRepository, ProductRepository},
};

#[derive(Clone)]
pub struct ProductGrpcService {
    repository: Arc<PostgresProductRepository>,
}

impl ProductGrpcService {
    pub fn new(repository: Arc<PostgresProductRepository>) -> Self {
        Self { repository }
    }
}

#[tonic::async_trait]
impl ProductService for ProductGrpcService {
    async fn get_product_by_id(
        &self,
        request: Request<GetProductByIdRequest>,
    ) -> Result<Response<ProductResponse>, Status> {
        let request = request.into_inner();
        let product = self
            .repository
            .get_by_id(request.id)
            .await
            .map_err(internal_error)?;

        let response = product.unwrap_or(Product {
            id: request.id,
            name: format!("Product {}", request.id),
            description: "Product not found in sample list".to_string(),
            price: Decimal::ZERO,
            stock_quantity: 0,
            created_at: Utc::now(),
            updated_at: None,
        });

        Ok(Response::new(map_product(response)))
    }

    async fn get_all_products(
        &self,
        _request: Request<GetAllProductsRequest>,
    ) -> Result<Response<GetAllProductsResponse>, Status> {
        let products = self.repository.get_all().await.map_err(internal_error)?;
        let response = GetAllProductsResponse {
            products: products.into_iter().map(map_product).collect(),
        };

        Ok(Response::new(response))
    }

    async fn create_product(
        &self,
        request: Request<CreateProductRequest>,
    ) -> Result<Response<ProductResponse>, Status> {
        let request = request.into_inner();
        let price = Decimal::from_f64(request.price)
            .ok_or_else(|| Status::invalid_argument("invalid decimal price"))?;

        let product = self
            .repository
            .create(CreateProductInput {
                name: request.name,
                description: request.description,
                price,
                stock_quantity: request.stock_quantity,
            })
            .await
            .map_err(internal_error)?;

        Ok(Response::new(map_product(product)))
    }

    async fn update_product(
        &self,
        request: Request<UpdateProductRequest>,
    ) -> Result<Response<ProductResponse>, Status> {
        let request = request.into_inner();
        let price = Decimal::from_f64(request.price)
            .ok_or_else(|| Status::invalid_argument("invalid decimal price"))?;

        let product = self
            .repository
            .update(
                request.id,
                UpdateProductInput {
                    name: request.name,
                    description: request.description,
                    price,
                    stock_quantity: request.stock_quantity,
                },
            )
            .await
            .map_err(internal_error)?;

        Ok(Response::new(map_product(product)))
    }

    async fn delete_product(
        &self,
        request: Request<DeleteProductRequest>,
    ) -> Result<Response<DeleteProductResponse>, Status> {
        let deleted = self
            .repository
            .delete(request.into_inner().id)
            .await
            .map_err(internal_error)?;

        Ok(Response::new(DeleteProductResponse { deleted }))
    }

    async fn delete_all_products(
        &self,
        _request: Request<DeleteAllProductsRequest>,
    ) -> Result<Response<DeleteAllProductsResponse>, Status> {
        let deleted = self.repository.delete_all().await.map_err(internal_error)?;
        Ok(Response::new(DeleteAllProductsResponse { deleted }))
    }
}

fn internal_error(error: sqlx::Error) -> Status {
    Status::internal(error.to_string())
}

fn map_product(product: Product) -> ProductResponse {
    ProductResponse {
        id: product.id,
        name: product.name,
        description: product.description,
        price: product.price.to_f64().unwrap_or_default(),
        stock_quantity: product.stock_quantity,
        created_at: product.created_at.to_rfc3339(),
        updated_at: product
            .updated_at
            .map(|value| value.to_rfc3339())
            .unwrap_or_default(),
    }
}
