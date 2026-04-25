use chrono::{DateTime, Utc};
use rust_decimal::Decimal;

#[derive(Debug, Clone, sqlx::FromRow)]
pub struct Product {
    #[sqlx(rename = "Id")]
    pub id: i32,
    #[sqlx(rename = "Name")]
    pub name: String,
    #[sqlx(rename = "Description")]
    pub description: String,
    #[sqlx(rename = "Category")]
    pub category: String,
    #[sqlx(rename = "Images")]
    pub images: Vec<String>,
    #[sqlx(rename = "Price")]
    pub price: Decimal,
    #[sqlx(rename = "StockQuantity")]
    pub stock_quantity: i32,
    #[sqlx(rename = "CreatedAt")]
    pub created_at: DateTime<Utc>,
    #[sqlx(rename = "UpdatedAt")]
    pub updated_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone)]
pub struct CreateProductInput {
    pub name: String,
    pub description: String,
    pub category: String,
    pub images: Vec<String>,
    pub price: Decimal,
    pub stock_quantity: i32,
}

#[derive(Debug, Clone)]
pub struct UpdateProductInput {
    pub name: String,
    pub description: String,
    pub category: String,
    pub images: Vec<String>,
    pub price: Decimal,
    pub stock_quantity: i32,
}
