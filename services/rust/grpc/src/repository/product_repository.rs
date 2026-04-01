use async_trait::async_trait;
use chrono::Utc;
use sqlx::PgPool;

use crate::product::{CreateProductInput, Product, UpdateProductInput};

#[async_trait]
pub trait ProductRepository: Send + Sync {
    async fn get_by_id(&self, id: i32) -> Result<Option<Product>, sqlx::Error>;
    async fn get_all(&self) -> Result<Vec<Product>, sqlx::Error>;
    async fn create(&self, request: CreateProductInput) -> Result<Product, sqlx::Error>;
    async fn update(&self, id: i32, request: UpdateProductInput) -> Result<Product, sqlx::Error>;
    async fn delete(&self, id: i32) -> Result<bool, sqlx::Error>;
    async fn delete_all(&self) -> Result<bool, sqlx::Error>;
}

#[derive(Clone)]
pub struct PostgresProductRepository {
    pool: PgPool,
}

impl PostgresProductRepository {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }
}

#[async_trait]
impl ProductRepository for PostgresProductRepository {
    async fn get_by_id(&self, id: i32) -> Result<Option<Product>, sqlx::Error> {
        sqlx::query_as::<_, Product>(
            r#"
            SELECT "Id", "Name", "Description", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
            FROM public.products
            WHERE "Id" = $1
            "#,
        )
        .bind(id)
        .fetch_optional(&self.pool)
        .await
    }

    async fn get_all(&self) -> Result<Vec<Product>, sqlx::Error> {
        sqlx::query_as::<_, Product>(
            r#"
            SELECT "Id", "Name", "Description", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
            FROM public.products
            ORDER BY "Id"
            "#,
        )
        .fetch_all(&self.pool)
        .await
    }

    async fn create(&self, request: CreateProductInput) -> Result<Product, sqlx::Error> {
        sqlx::query_as::<_, Product>(
            r#"
            INSERT INTO public.products ("Name", "Description", "Price", "StockQuantity", "CreatedAt")
            VALUES ($1, $2, $3, $4, $5)
            RETURNING "Id", "Name", "Description", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
            "#,
        )
        .bind(request.name)
        .bind(request.description)
        .bind(request.price)
        .bind(request.stock_quantity)
        .bind(Utc::now())
        .fetch_one(&self.pool)
        .await
    }

    async fn update(&self, id: i32, request: UpdateProductInput) -> Result<Product, sqlx::Error> {
        let existing = self.get_by_id(id).await?;

        if existing.is_none() {
            return self
                .create(CreateProductInput {
                    name: request.name,
                    description: request.description,
                    price: request.price,
                    stock_quantity: request.stock_quantity,
                })
                .await;
        }

        sqlx::query_as::<_, Product>(
            r#"
            UPDATE public.products
            SET "Name" = $1,
                "Description" = $2,
                "Price" = $3,
                "StockQuantity" = $4,
                "UpdatedAt" = $5
            WHERE "Id" = $6
            RETURNING "Id", "Name", "Description", "Price", "StockQuantity", "CreatedAt", "UpdatedAt"
            "#,
        )
        .bind(request.name)
        .bind(request.description)
        .bind(request.price)
        .bind(request.stock_quantity)
        .bind(Utc::now())
        .bind(id)
        .fetch_one(&self.pool)
        .await
    }

    async fn delete(&self, id: i32) -> Result<bool, sqlx::Error> {
        let result = sqlx::query(r#"DELETE FROM public.products WHERE "Id" = $1"#)
            .bind(id)
            .execute(&self.pool)
            .await?;

        Ok(result.rows_affected() > 0)
    }

    async fn delete_all(&self) -> Result<bool, sqlx::Error> {
        let result = sqlx::query(r#"DELETE FROM public.products"#)
            .execute(&self.pool)
            .await?;

        Ok(result.rows_affected() > 0)
    }
}
