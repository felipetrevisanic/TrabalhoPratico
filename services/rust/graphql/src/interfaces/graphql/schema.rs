use std::sync::Arc;

use async_graphql::{Context, EmptySubscription, Object, Result, Schema};
use chrono::Utc;
use rust_decimal::Decimal;

use crate::{
    application::interfaces::product_service::ProductService,
    domain::entities::product::{CreateProductInput, Product, UpdateProductInput},
};

pub type ProductSchema = Schema<QueryRoot, MutationRoot, EmptySubscription>;

pub struct QueryRoot;

#[Object]
impl QueryRoot {
    async fn product_by_id(&self, context: &Context<'_>, id: i32) -> Result<Product> {
        let service = context.data::<Arc<dyn ProductService>>()?;
        let product = service.get_by_id(id).await?;

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
        let service = context.data::<Arc<dyn ProductService>>()?;
        Ok(service.get_all().await?)
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
        let service = context.data::<Arc<dyn ProductService>>()?;
        Ok(service.create(input).await?)
    }

    async fn update_product(
        &self,
        context: &Context<'_>,
        id: i32,
        input: UpdateProductInput,
    ) -> Result<Product> {
        let service = context.data::<Arc<dyn ProductService>>()?;
        Ok(service.update(id, input).await?)
    }

    async fn delete_product(&self, context: &Context<'_>, id: i32) -> Result<bool> {
        let service = context.data::<Arc<dyn ProductService>>()?;
        Ok(service.delete(id).await?)
    }

    async fn delete_all_products(&self, context: &Context<'_>) -> Result<bool> {
        let service = context.data::<Arc<dyn ProductService>>()?;
        Ok(service.delete_all().await?)
    }
}
