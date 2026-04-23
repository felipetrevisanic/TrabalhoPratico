use async_trait::async_trait;

use crate::domain::entities::product::{CreateProductInput, Product, UpdateProductInput};

#[async_trait]
pub trait ProductRepository: Send + Sync {
    async fn get_by_id(&self, id: i32) -> Result<Option<Product>, sqlx::Error>;
    async fn get_all(&self) -> Result<Vec<Product>, sqlx::Error>;
    async fn create(&self, request: CreateProductInput) -> Result<Product, sqlx::Error>;
    async fn update(&self, id: i32, request: UpdateProductInput) -> Result<Product, sqlx::Error>;
    async fn delete(&self, id: i32) -> Result<bool, sqlx::Error>;
    async fn delete_all(&self) -> Result<bool, sqlx::Error>;
}
