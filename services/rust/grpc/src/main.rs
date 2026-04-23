mod application;
mod config;
mod domain;
mod infrastructure;
mod interfaces;

use std::sync::Arc;

use sqlx::postgres::PgPoolOptions;
use tonic::transport::Server;

use crate::application::{
    interfaces::product_service::ProductService, service::product_service::ProductServiceImpl,
};
use crate::config::AppConfig;
use crate::infrastructure::repositories::product_repository::PostgresProductRepository;
use crate::interfaces::grpc::{
    ProductGrpcService, productv1::product_service_server::ProductServiceServer,
};

#[tokio::main]
async fn main() {
    let config = AppConfig::from_env();

    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(&config.database_url())
        .await
        .expect("failed to connect to PostgreSQL");

    let repository = Arc::new(PostgresProductRepository::new(pool));
    let application_service: Arc<dyn ProductService> =
        Arc::new(ProductServiceImpl::new(repository));
    let service = ProductGrpcService::new(application_service);

    let address = config.grpc_address().parse().expect("invalid gRPC address");

    println!("gRPC API listening on {address}");

    Server::builder()
        .add_service(ProductServiceServer::new(service))
        .serve(address)
        .await
        .expect("failed to serve gRPC server");
}
