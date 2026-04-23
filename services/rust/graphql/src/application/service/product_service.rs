use std::sync::Arc;

use async_trait::async_trait;

use crate::{
    application::interfaces::product_service::ProductService,
    domain::{
        entities::product::{CreateProductInput, Product, UpdateProductInput},
        interfaces::product_repository::ProductRepository,
    },
};

#[derive(Clone)]
pub struct ProductServiceImpl {
    repository: Arc<dyn ProductRepository>,
}

impl ProductServiceImpl {
    pub fn new(repository: Arc<dyn ProductRepository>) -> Self {
        Self { repository }
    }
}

#[async_trait]
impl ProductService for ProductServiceImpl {
    async fn get_by_id(&self, id: i32) -> Result<Option<Product>, sqlx::Error> {
        self.repository.get_by_id(id).await
    }

    async fn get_all(&self) -> Result<Vec<Product>, sqlx::Error> {
        self.repository.get_all().await
    }

    async fn create(&self, request: CreateProductInput) -> Result<Product, sqlx::Error> {
        self.repository.create(request).await
    }

    async fn update(&self, id: i32, request: UpdateProductInput) -> Result<Product, sqlx::Error> {
        self.repository.update(id, request).await
    }

    async fn delete(&self, id: i32) -> Result<bool, sqlx::Error> {
        self.repository.delete(id).await
    }

    async fn delete_all(&self) -> Result<bool, sqlx::Error> {
        self.repository.delete_all().await
    }
}
