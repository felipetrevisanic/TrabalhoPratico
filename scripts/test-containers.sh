#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GRPCURL_BIN="${ROOT_DIR}/.tools/bin/grpcurl"
PROTO_FILE="${ROOT_DIR}/shared/proto/product/v1/product.proto"

wait_http() {
  local url="$1"
  local name="$2"

  for _ in $(seq 1 60); do
    if curl -fsS "$url" >/dev/null 2>&1; then
      return 0
    fi
    sleep 2
  done

  echo "Servico HTTP indisponivel: ${name} (${url})" >&2
  return 1
}

wait_graphql() {
  local url="$1"
  local name="$2"

  for _ in $(seq 1 60); do
    if curl -fsS -X POST "$url" \
      -H 'Content-Type: application/json' \
      -d '{"query":"{ allProducts { id name } }"}' >/dev/null 2>&1; then
      return 0
    fi
    sleep 2
  done

  echo "Servico GraphQL indisponivel: ${name} (${url})" >&2
  return 1
}

wait_grpc() {
  local address="$1"
  local name="$2"

  for _ in $(seq 1 60); do
    if "${GRPCURL_BIN}" -plaintext -import-path "${ROOT_DIR}/shared/proto" -proto "${PROTO_FILE}" \
      "$address" product.v1.ProductService/GetAllProducts >/dev/null 2>&1; then
      return 0
    fi
    sleep 2
  done

  echo "Servico gRPC indisponivel: ${name} (${address})" >&2
  return 1
}

assert_contains() {
  local haystack="$1"
  local needle="$2"
  local message="$3"
  local normalized_haystack
  local normalized_needle

  normalized_haystack="$(printf '%s' "$haystack" | tr -d '[:space:]')"
  normalized_needle="$(printf '%s' "$needle" | tr -d '[:space:]')"

  if [[ "$normalized_haystack" != *"$normalized_needle"* ]]; then
    echo "Falha: ${message}" >&2
    echo "Resposta recebida: ${haystack}" >&2
    exit 1
  fi
}

wait_http "http://localhost:5107/Product/all" "csharp-rest"
wait_http "http://localhost:8080/Product/all" "go-rest"
wait_http "http://localhost:8083/Product/all" "rust-rest"
wait_graphql "http://localhost:5207/graphql" "csharp-graphql"
wait_graphql "http://localhost:8090/graphql" "go-graphql"
wait_graphql "http://localhost:8091/graphql" "rust-graphql"
wait_grpc "localhost:5307" "csharp-grpc"
wait_grpc "localhost:50051" "go-grpc"
wait_grpc "localhost:50052" "rust-grpc"

csharp_rest_all="$(curl -fsS http://localhost:5107/Product/all)"
go_rest_all="$(curl -fsS http://localhost:8080/Product/all)"
rust_rest_all="$(curl -fsS http://localhost:8083/Product/all)"

assert_contains "$csharp_rest_all" '"name":"Notebook"' "C# REST deve retornar o produto seedado Notebook"
assert_contains "$csharp_rest_all" '"name":"Mouse"' "C# REST deve retornar o produto seedado Mouse"
assert_contains "$go_rest_all" '"name":"Notebook"' "Go REST deve retornar o produto seedado Notebook"
assert_contains "$go_rest_all" '"name":"Mouse"' "Go REST deve retornar o produto seedado Mouse"
assert_contains "$rust_rest_all" '"name":"Notebook"' "Rust REST deve retornar o produto seedado Notebook"
assert_contains "$rust_rest_all" '"name":"Mouse"' "Rust REST deve retornar o produto seedado Mouse"

csharp_rest_missing="$(curl -fsS 'http://localhost:5107/Product?id=999')"
go_rest_missing="$(curl -fsS 'http://localhost:8080/Product?id=999')"
rust_rest_missing="$(curl -fsS 'http://localhost:8083/Product?id=999')"

assert_contains "$csharp_rest_missing" '"name":"Product 999"' "C# REST deve manter o placeholder para item ausente"
assert_contains "$go_rest_missing" '"name":"Product 999"' "Go REST deve manter o placeholder para item ausente"
assert_contains "$rust_rest_missing" '"name":"Product 999"' "Rust REST deve manter o placeholder para item ausente"

csharp_graphql_all="$(curl -fsS -X POST http://localhost:5207/graphql -H 'Content-Type: application/json' -d '{"query":"{ allProducts { id name } }"}')"
go_graphql_all="$(curl -fsS -X POST http://localhost:8090/graphql -H 'Content-Type: application/json' -d '{"query":"{ allProducts { id name } }"}')"
rust_graphql_all="$(curl -fsS -X POST http://localhost:8091/graphql -H 'Content-Type: application/json' -d '{"query":"{ allProducts { id name } }"}')"

assert_contains "$csharp_graphql_all" '"name":"Notebook"' "C# GraphQL deve retornar Notebook"
assert_contains "$go_graphql_all" '"name":"Notebook"' "Go GraphQL deve retornar Notebook"
assert_contains "$rust_graphql_all" '"name":"Notebook"' "Rust GraphQL deve retornar Notebook"

csharp_graphql_missing="$(curl -fsS -X POST http://localhost:5207/graphql -H 'Content-Type: application/json' -d '{"query":"{ productById(id: 999) { id name description } }"}')"
go_graphql_missing="$(curl -fsS -X POST http://localhost:8090/graphql -H 'Content-Type: application/json' -d '{"query":"{ productById(id: 999) { id name description } }"}')"
rust_graphql_missing="$(curl -fsS -X POST http://localhost:8091/graphql -H 'Content-Type: application/json' -d '{"query":"{ productById(id: 999) { id name description } }"}')"

assert_contains "$csharp_graphql_missing" '"name":"Product 999"' "C# GraphQL deve manter o placeholder para item ausente"
assert_contains "$go_graphql_missing" '"name":"Product 999"' "Go GraphQL deve manter o placeholder para item ausente"
assert_contains "$rust_graphql_missing" '"name":"Product 999"' "Rust GraphQL deve manter o placeholder para item ausente"

csharp_grpc_all="$("${GRPCURL_BIN}" -plaintext -import-path "${ROOT_DIR}/shared/proto" -proto "${PROTO_FILE}" \
  localhost:5307 product.v1.ProductService/GetAllProducts)"
go_grpc_all="$("${GRPCURL_BIN}" -plaintext -import-path "${ROOT_DIR}/shared/proto" -proto "${PROTO_FILE}" \
  localhost:50051 product.v1.ProductService/GetAllProducts)"
rust_grpc_all="$("${GRPCURL_BIN}" -plaintext -import-path "${ROOT_DIR}/shared/proto" -proto "${PROTO_FILE}" \
  localhost:50052 product.v1.ProductService/GetAllProducts)"

assert_contains "$csharp_grpc_all" '"name": "Notebook"' "C# gRPC deve retornar Notebook"
assert_contains "$go_grpc_all" '"name": "Notebook"' "Go gRPC deve retornar Notebook"
assert_contains "$rust_grpc_all" '"name": "Notebook"' "Rust gRPC deve retornar Notebook"

csharp_grpc_missing="$("${GRPCURL_BIN}" -plaintext -import-path "${ROOT_DIR}/shared/proto" -proto "${PROTO_FILE}" \
  -d '{"id":999}' localhost:5307 product.v1.ProductService/GetProductById)"
go_grpc_missing="$("${GRPCURL_BIN}" -plaintext -import-path "${ROOT_DIR}/shared/proto" -proto "${PROTO_FILE}" \
  -d '{"id":999}' localhost:50051 product.v1.ProductService/GetProductById)"
rust_grpc_missing="$("${GRPCURL_BIN}" -plaintext -import-path "${ROOT_DIR}/shared/proto" -proto "${PROTO_FILE}" \
  -d '{"id":999}' localhost:50052 product.v1.ProductService/GetProductById)"

assert_contains "$csharp_grpc_missing" '"name": "Product 999"' "C# gRPC deve manter o placeholder para item ausente"
assert_contains "$go_grpc_missing" '"name": "Product 999"' "Go gRPC deve manter o placeholder para item ausente"
assert_contains "$rust_grpc_missing" '"name": "Product 999"' "Rust gRPC deve manter o placeholder para item ausente"

echo "Todos os containers responderam conforme esperado."
