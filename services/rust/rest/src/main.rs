mod application;
mod config;
mod domain;
mod infrastructure;
mod interfaces;

use std::net::SocketAddr;
use std::sync::Arc;

use axum::Router;
use sqlx::postgres::PgPoolOptions;

use crate::application::service::product_service::ProductServiceImpl;
use crate::config::AppConfig;
use crate::infrastructure::repositories::product_repository::PostgresProductRepository;
use crate::interfaces::rest::create_router;

#[tokio::main]
async fn main() {
    let config = AppConfig::from_env();

    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(&config.database_url())
        .await
        .expect("failed to connect to PostgreSQL");

    let repository = Arc::new(PostgresProductRepository::new(pool));
    let service = Arc::new(ProductServiceImpl::new(repository));
    let app: Router = create_router(service);

    let listener = tokio::net::TcpListener::bind(config.http_address())
        .await
        .expect("failed to bind TCP listener");

    let address: SocketAddr = listener
        .local_addr()
        .expect("failed to read bound listener address");

    println!("REST API listening on {address}");

    axum::serve(listener, app)
        .await
        .expect("failed to serve HTTP application");
}
