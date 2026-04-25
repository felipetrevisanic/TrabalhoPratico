use std::sync::Arc;

use crate::{
    application::interfaces::product_service::ProductService,
    domain::entities::product::{CreateProductInput, Product, UpdateProductInput},
};
use async_graphql::{Context, EmptySubscription, Object, Result, Schema};

pub type ProductSchema = Schema<QueryRoot, MutationRoot, EmptySubscription>;

pub struct QueryRoot;

#[Object]
impl QueryRoot {
    async fn product_by_id(&self, context: &Context<'_>, id: i32) -> Result<Option<Product>> {
        let service = context.data::<Arc<dyn ProductService>>()?;
        Ok(service.get_by_id(id).await?)
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
