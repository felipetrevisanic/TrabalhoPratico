#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GRPCURL_BIN="${ROOT_DIR}/.tools/bin/grpcurl"
PROTO_IMPORT_PATH="${ROOT_DIR}/shared/proto"
PROTO_FILE="product/v1/product.proto"

REST_CREATE='{"name":"Produto REST","description":"Criado via teste REST","category":"Hardware","images":["https://img.local/rest-1.png","https://img.local/rest-2.png"],"price":49.90,"stockQuantity":7}'
REST_UPDATE='{"name":"Produto REST Atualizado","description":"Atualizado via teste REST","category":"Hardware Updated","images":["https://img.local/rest-1.png?v=2","https://img.local/rest-2.png?v=2"],"price":54.15,"stockQuantity":11}'
GRAPHQL_CREATE='{"name":"Produto GraphQL","description":"Criado via teste GraphQL","category":"Peripherals","images":["https://img.local/graphql-1.png","https://img.local/graphql-2.png"],"price":63.40,"stockQuantity":9}'
GRAPHQL_UPDATE='{"name":"Produto GraphQL Atualizado","description":"Atualizado via teste GraphQL","category":"Peripherals Updated","images":["https://img.local/graphql-1.png?v=2","https://img.local/graphql-2.png?v=2"],"price":68.55,"stockQuantity":13}'
GRPC_CREATE='{"name":"Produto gRPC","description":"Criado via teste gRPC","category":"Services","images":["https://img.local/grpc-1.png","https://img.local/grpc-2.png"],"price":71.25,"stock_quantity":5}'
GRPC_UPDATE='{"name":"Produto gRPC Atualizado","description":"Atualizado via teste gRPC","category":"Services Updated","images":["https://img.local/grpc-1.png?v=2","https://img.local/grpc-2.png?v=2"],"price":74.50,"stock_quantity":8}'

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
    if "${GRPCURL_BIN}" -plaintext -import-path "${PROTO_IMPORT_PATH}" -proto "${PROTO_FILE}" \
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

extract_id() {
  local body="$1"

  printf '%s' "$body" \
    | tr -d '\n' \
    | sed -E 's/.*"id"[[:space:]]*:[[:space:]]*([0-9]+).*/\1/'
}

assert_http_status() {
  local url="$1"
  local expected="$2"
  local message="$3"
  local method="${4:-GET}"
  local body_file="/tmp/codex-http-body.$$"

  local status
  status="$(curl -sS -o "$body_file" -w '%{http_code}' -X "$method" "$url")"

  if [[ "$status" != "$expected" ]]; then
    echo "Falha: ${message}" >&2
    echo "Status recebido: ${status}" >&2
    echo "Body recebido: $(cat "$body_file")" >&2
    rm -f "$body_file"
    exit 1
  fi

  rm -f "$body_file"
}

test_rest_api() {
  local name="$1"
  local base_url="$2"

  curl -fsS -X DELETE "${base_url}/Product/all" >/dev/null 2>&1 || true

  local created
  created="$(curl -fsS -X POST "${base_url}/Product" -H 'Content-Type: application/json' -d "${REST_CREATE}")"
  assert_contains "$created" '"name":"Produto REST"' "${name}: create deve retornar name"
  assert_contains "$created" '"category":"Hardware"' "${name}: create deve retornar category"
  assert_contains "$created" '"images":["https://img.local/rest-1.png","https://img.local/rest-2.png"]' "${name}: create deve retornar images"
  local id
  id="$(extract_id "$created")"

  local by_id
  by_id="$(curl -fsS "${base_url}/Product?id=${id}")"
  assert_contains "$by_id" '"description":"Criado via teste REST"' "${name}: get by id deve retornar description"
  assert_contains "$by_id" '"stockQuantity":7' "${name}: get by id deve retornar stockQuantity"

  local all_products
  all_products="$(curl -fsS "${base_url}/Product/all")"
  assert_contains "$all_products" '"name":"Produto REST"' "${name}: get all deve listar produto criado"
  assert_contains "$all_products" '"category":"Hardware"' "${name}: get all deve listar category"

  local updated
  updated="$(curl -fsS -X PUT "${base_url}/Product/${id}" -H 'Content-Type: application/json' -d "${REST_UPDATE}")"
  assert_contains "$updated" '"name":"Produto REST Atualizado"' "${name}: update deve retornar novo name"
  assert_contains "$updated" '"category":"Hardware Updated"' "${name}: update deve retornar nova category"
  assert_contains "$updated" '"images":["https://img.local/rest-1.png?v=2","https://img.local/rest-2.png?v=2"]' "${name}: update deve retornar novas images"
  assert_contains "$updated" '"stockQuantity":11' "${name}: update deve retornar novo stockQuantity"

  assert_http_status "${base_url}/Product?id=999999" "404" "${name}: get by id ausente deve retornar 404"
}

test_graphql_api() {
  local name="$1"
  local url="$2"

  curl -fsS -X POST "$url" \
    -H 'Content-Type: application/json' \
    -d '{"query":"mutation { deleteAllProducts }"}' >/dev/null

  local create_payload
  create_payload="$(cat <<EOF
{"query":"mutation CreateProduct(\$input: CreateProductInput!) { createProduct(input: \$input) { id name description category images price stockQuantity } }","variables":{"input":${GRAPHQL_CREATE}}}
EOF
)"

  local created
  created="$(curl -fsS -X POST "$url" -H 'Content-Type: application/json' -d "$create_payload")"
  assert_contains "$created" '"name":"Produto GraphQL"' "${name}: create deve retornar name"
  assert_contains "$created" '"category":"Peripherals"' "${name}: create deve retornar category"
  assert_contains "$created" '"images":["https://img.local/graphql-1.png","https://img.local/graphql-2.png"]' "${name}: create deve retornar images"
  local id
  id="$(extract_id "$created")"

  local by_id_payload
  by_id_payload="$(cat <<EOF
{"query":"query ProductById(\$id: Int!) { productById(id: \$id) { id name description category images price stockQuantity } }","variables":{"id":${id}}}
EOF
)"

  local by_id
  by_id="$(curl -fsS -X POST "$url" -H 'Content-Type: application/json' -d "$by_id_payload")"
  assert_contains "$by_id" '"description":"Criado via teste GraphQL"' "${name}: get by id deve retornar description"
  assert_contains "$by_id" '"stockQuantity":9' "${name}: get by id deve retornar stockQuantity"

  local all_products
  all_products="$(curl -fsS -X POST "$url" -H 'Content-Type: application/json' -d '{"query":"{ allProducts { id name description category images price stockQuantity } }"}')"
  assert_contains "$all_products" '"name":"Produto GraphQL"' "${name}: get all deve listar produto criado"
  assert_contains "$all_products" '"category":"Peripherals"' "${name}: get all deve listar category"

  local update_payload
  update_payload="$(cat <<EOF
{"query":"mutation UpdateProduct(\$id: Int!, \$input: UpdateProductInput!) { updateProduct(id: \$id, input: \$input) { id name description category images price stockQuantity } }","variables":{"id":${id},"input":${GRAPHQL_UPDATE}}}
EOF
)"

  local updated
  updated="$(curl -fsS -X POST "$url" -H 'Content-Type: application/json' -d "$update_payload")"
  assert_contains "$updated" '"name":"Produto GraphQL Atualizado"' "${name}: update deve retornar novo name"
  assert_contains "$updated" '"category":"Peripherals Updated"' "${name}: update deve retornar nova category"
  assert_contains "$updated" '"images":["https://img.local/graphql-1.png?v=2","https://img.local/graphql-2.png?v=2"]' "${name}: update deve retornar novas images"
  assert_contains "$updated" '"stockQuantity":13' "${name}: update deve retornar novo stockQuantity"

  local missing
  missing="$(curl -fsS -X POST "$url" -H 'Content-Type: application/json' -d '{"query":"{ productById(id: 999999) { id name } }"}')"
  assert_contains "$missing" '"productById":null' "${name}: item ausente deve retornar null"
}

test_grpc_api() {
  local name="$1"
  local address="$2"

  "${GRPCURL_BIN}" -plaintext -import-path "${PROTO_IMPORT_PATH}" -proto "${PROTO_FILE}" \
    "$address" product.v1.ProductService/DeleteAllProducts >/dev/null

  local created
  created="$("${GRPCURL_BIN}" -plaintext -import-path "${PROTO_IMPORT_PATH}" -proto "${PROTO_FILE}" \
    -d "${GRPC_CREATE}" "$address" product.v1.ProductService/CreateProduct)"
  assert_contains "$created" '"name": "Produto gRPC"' "${name}: create deve retornar name"
  assert_contains "$created" '"category": "Services"' "${name}: create deve retornar category"
  assert_contains "$created" '"images": [' "${name}: create deve retornar images"
  local id
  id="$(extract_id "$created")"

  local by_id
  by_id="$("${GRPCURL_BIN}" -plaintext -import-path "${PROTO_IMPORT_PATH}" -proto "${PROTO_FILE}" \
    -d "{\"id\":${id}}" "$address" product.v1.ProductService/GetProductById)"
  assert_contains "$by_id" '"description": "Criado via teste gRPC"' "${name}: get by id deve retornar description"
  assert_contains "$by_id" '"stockQuantity": 5' "${name}: get by id deve retornar stockQuantity"

  local all_products
  all_products="$("${GRPCURL_BIN}" -plaintext -import-path "${PROTO_IMPORT_PATH}" -proto "${PROTO_FILE}" \
    "$address" product.v1.ProductService/GetAllProducts)"
  assert_contains "$all_products" '"name": "Produto gRPC"' "${name}: get all deve listar produto criado"
  assert_contains "$all_products" '"category": "Services"' "${name}: get all deve listar category"

  local updated
  updated="$("${GRPCURL_BIN}" -plaintext -import-path "${PROTO_IMPORT_PATH}" -proto "${PROTO_FILE}" \
    -d "{\"id\":${id},\"name\":\"Produto gRPC Atualizado\",\"description\":\"Atualizado via teste gRPC\",\"category\":\"Services Updated\",\"images\":[\"https://img.local/grpc-1.png?v=2\",\"https://img.local/grpc-2.png?v=2\"],\"price\":74.50,\"stock_quantity\":8}" \
    "$address" product.v1.ProductService/UpdateProduct)"
  assert_contains "$updated" '"name": "Produto gRPC Atualizado"' "${name}: update deve retornar novo name"
  assert_contains "$updated" '"category": "Services Updated"' "${name}: update deve retornar nova category"
  assert_contains "$updated" '"stockQuantity": 8' "${name}: update deve retornar novo stockQuantity"

  local missing_output
  local missing_status=0
  missing_output="$("${GRPCURL_BIN}" -plaintext -import-path "${PROTO_IMPORT_PATH}" -proto "${PROTO_FILE}" \
    -d '{"id":999999}' "$address" product.v1.ProductService/GetProductById 2>&1)" || missing_status=$?

  if [[ "$missing_status" -eq 0 ]]; then
    echo "Falha: ${name}: item ausente via gRPC deveria retornar erro" >&2
    exit 1
  fi

  assert_contains "$missing_output" 'NotFound' "${name}: item ausente via gRPC deve retornar NotFound"
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

test_rest_api "csharp-rest" "http://localhost:5107"
test_rest_api "go-rest" "http://localhost:8080"
test_rest_api "rust-rest" "http://localhost:8083"

test_graphql_api "csharp-graphql" "http://localhost:5207/graphql"
test_graphql_api "go-graphql" "http://localhost:8090/graphql"
test_graphql_api "rust-graphql" "http://localhost:8091/graphql"

test_grpc_api "csharp-grpc" "localhost:5307"
test_grpc_api "go-grpc" "localhost:50051"
test_grpc_api "rust-grpc" "localhost:50052"

echo "Todos os containers responderam conforme esperado."
