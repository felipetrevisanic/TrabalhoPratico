mod application;
mod config;
mod domain;
mod infrastructure;
mod interfaces;

use std::net::SocketAddr;
use std::sync::Arc;

use async_graphql::Schema;
use async_graphql_axum::{GraphQLRequest, GraphQLResponse};
use axum::{
    Router,
    extract::State,
    response::{Html, IntoResponse},
    routing::get,
};
use sqlx::postgres::PgPoolOptions;

use crate::application::{
    interfaces::product_service::ProductService, service::product_service::ProductServiceImpl,
};
use crate::config::AppConfig;
use crate::infrastructure::repositories::product_repository::PostgresProductRepository;
use crate::interfaces::graphql::{MutationRoot, ProductSchema, QueryRoot};

#[tokio::main]
async fn main() {
    let config = AppConfig::from_env();

    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(&config.database_url())
        .await
        .expect("failed to connect to PostgreSQL");

    let repository = Arc::new(PostgresProductRepository::new(pool));
    let service: Arc<dyn ProductService> = Arc::new(ProductServiceImpl::new(repository));
    let schema = Schema::build(QueryRoot, MutationRoot, async_graphql::EmptySubscription)
        .data(service)
        .finish();

    let app = Router::new()
        .route("/graphql", get(graphql_playground).post(graphql_handler))
        .with_state(schema);

    let listener = tokio::net::TcpListener::bind(config.http_address())
        .await
        .expect("failed to bind TCP listener");

    let address: SocketAddr = listener
        .local_addr()
        .expect("failed to read bound listener address");

    println!("GraphQL API listening on {address}");

    axum::serve(listener, app)
        .await
        .expect("failed to serve GraphQL application");
}

async fn graphql_handler(
    State(schema): State<ProductSchema>,
    request: GraphQLRequest,
) -> GraphQLResponse {
    schema.execute(request.into_inner()).await.into()
}

async fn graphql_playground() -> impl IntoResponse {
    Html(async_graphql::http::playground_source(
        async_graphql::http::GraphQLPlaygroundConfig::new("/graphql"),
    ))
}
