use std::sync::Arc;

use axum::{
    Json,
    extract::{Path, Query, State},
    http::StatusCode,
};
use serde::Deserialize;

use crate::{
    application::interfaces::product_service::ProductService,
    domain::entities::product::{CreateProductRequest, Product, UpdateProductRequest},
};

#[derive(Deserialize)]
pub struct GetProductQuery {
    id: i32,
}

pub async fn get_product_by_id(
    State(service): State<Arc<dyn ProductService>>,
    Query(query): Query<GetProductQuery>,
) -> Result<Json<Product>, StatusCode> {
    let product = service
        .get_by_id(query.id)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let response = product.ok_or(StatusCode::NOT_FOUND)?;

    Ok(Json(response))
}

pub async fn get_all_products(
    State(service): State<Arc<dyn ProductService>>,
) -> Result<Json<Vec<Product>>, StatusCode> {
    let products = service
        .get_all()
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(products))
}

pub async fn create_product(
    State(service): State<Arc<dyn ProductService>>,
    Json(request): Json<CreateProductRequest>,
) -> Result<Json<Product>, StatusCode> {
    let product = service
        .create(request)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(product))
}

pub async fn update_product(
    State(service): State<Arc<dyn ProductService>>,
    Path(id): Path<i32>,
    Json(request): Json<UpdateProductRequest>,
) -> Result<Json<Product>, StatusCode> {
    let product = service
        .update(id, request)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    Ok(Json(product))
}

pub async fn delete_product(
    State(service): State<Arc<dyn ProductService>>,
    Path(id): Path<i32>,
) -> Result<StatusCode, StatusCode> {
    let deleted = service
        .delete(id)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    if !deleted {
        return Ok(StatusCode::NOT_FOUND);
    }

    Ok(StatusCode::NO_CONTENT)
}

pub async fn delete_all_products(
    State(service): State<Arc<dyn ProductService>>,
) -> Result<StatusCode, StatusCode> {
    let deleted = service
        .delete_all()
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    if !deleted {
        return Ok(StatusCode::NOT_FOUND);
    }

    Ok(StatusCode::NO_CONTENT)
}
