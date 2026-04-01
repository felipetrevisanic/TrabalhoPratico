mod config;
mod grpc;
mod product;
mod repository;

use std::sync::Arc;

use sqlx::postgres::PgPoolOptions;
use tonic::transport::Server;

use crate::config::AppConfig;
use crate::grpc::productv1::product_service_server::ProductServiceServer;
use crate::grpc::service::ProductGrpcService;
use crate::repository::product_repository::PostgresProductRepository;

#[tokio::main]
async fn main() {
    let config = AppConfig::from_env();

    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(&config.database_url())
        .await
        .expect("failed to connect to PostgreSQL");

    let repository = Arc::new(PostgresProductRepository::new(pool));
    let service = ProductGrpcService::new(repository);

    let address = config.grpc_address().parse().expect("invalid gRPC address");

    println!("gRPC API listening on {address}");

    Server::builder()
        .add_service(ProductServiceServer::new(service))
        .serve(address)
        .await
        .expect("failed to serve gRPC server");
}
