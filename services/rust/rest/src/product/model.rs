use chrono::{DateTime, Utc};
use rust_decimal::Decimal;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, sqlx::FromRow)]
pub struct Product {
    #[sqlx(rename = "Id")]
    pub id: i32,
    #[sqlx(rename = "Name")]
    pub name: String,
    #[sqlx(rename = "Description")]
    pub description: String,
    #[sqlx(rename = "Price")]
    pub price: Decimal,
    #[sqlx(rename = "StockQuantity")]
    #[serde(rename = "stockQuantity")]
    pub stock_quantity: i32,
    #[sqlx(rename = "CreatedAt")]
    #[serde(rename = "createdAt")]
    pub created_at: DateTime<Utc>,
    #[sqlx(rename = "UpdatedAt")]
    #[serde(rename = "updatedAt")]
    pub updated_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Deserialize)]
pub struct CreateProductRequest {
    pub name: String,
    pub description: String,
    pub price: Decimal,
    #[serde(rename = "stockQuantity")]
    pub stock_quantity: i32,
}

#[derive(Debug, Deserialize)]
pub struct UpdateProductRequest {
    pub name: String,
    pub description: String,
    pub price: Decimal,
    #[serde(rename = "stockQuantity")]
    pub stock_quantity: i32,
}
