mod handlers;

use std::sync::Arc;

use axum::{
    Router,
    routing::{get, put},
};

use crate::repository::product_repository::ProductRepository;

pub fn create_router(repository: Arc<dyn ProductRepository>) -> Router {
    Router::new()
        .route(
            "/Product",
            get(handlers::get_product_by_id).post(handlers::create_product),
        )
        .route(
            "/Product/all",
            get(handlers::get_all_products).delete(handlers::delete_all_products),
        )
        .route(
            "/Product/{id}",
            put(handlers::update_product).delete(handlers::delete_product),
        )
        .with_state(repository)
}
