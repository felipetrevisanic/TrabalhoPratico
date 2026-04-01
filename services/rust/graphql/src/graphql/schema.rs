use std::sync::Arc;

use async_graphql::{Context, EmptySubscription, Object, Result, Schema};
use chrono::Utc;
use rust_decimal::Decimal;

use crate::{
    product::{CreateProductInput, Product, UpdateProductInput},
    repository::product_repository::{PostgresProductRepository, ProductRepository},
};

pub type ProductSchema = Schema<QueryRoot, MutationRoot, EmptySubscription>;

pub struct QueryRoot;

#[Object]
impl QueryRoot {
    async fn product_by_id(&self, context: &Context<'_>, id: i32) -> Result<Product> {
        let repository = context.data::<Arc<PostgresProductRepository>>()?;
        let product = repository.get_by_id(id).await?;

        Ok(product.unwrap_or(Product {
            id,
            name: format!("Product {id}"),
            description: "Product not found in sample list".to_string(),
            price: Decimal::ZERO,
            stock_quantity: 0,
            created_at: Utc::now(),
            updated_at: None,
        }))
    }

    async fn all_products(&self, context: &Context<'_>) -> Result<Vec<Product>> {
        let repository = context.data::<Arc<PostgresProductRepository>>()?;
        Ok(repository.get_all().await?)
    }
}

pub struct MutationRoot;

#[Object]
impl MutationRoot {
    async fn create_product(
        &self,
        context: &Context<'_>,
        input: CreateProductInput,
    ) -> Result<Product> {
        let repository = context.data::<Arc<PostgresProductRepository>>()?;
        Ok(repository.create(input).await?)
    }

    async fn update_product(
        &self,
        context: &Context<'_>,
        id: i32,
        input: UpdateProductInput,
    ) -> Result<Product> {
        let repository = context.data::<Arc<PostgresProductRepository>>()?;
        Ok(repository.update(id, input).await?)
    }

    async fn delete_product(&self, context: &Context<'_>, id: i32) -> Result<bool> {
        let repository = context.data::<Arc<PostgresProductRepository>>()?;
        Ok(repository.delete(id).await?)
    }

    async fn delete_all_products(&self, context: &Context<'_>) -> Result<bool> {
        let repository = context.data::<Arc<PostgresProductRepository>>()?;
        Ok(repository.delete_all().await?)
    }
}
