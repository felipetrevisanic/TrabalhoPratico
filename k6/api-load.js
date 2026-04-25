import http from 'k6/http';
import grpc from 'k6/net/grpc';
import { check } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

import { resolveDataProfile, resolveTarget } from './config.js';
import { buildProducts, buildUpdatedProduct, pickIndexes } from './data.js';

const target = resolveTarget(__ENV.TARGET_ID);
const profile = resolveDataProfile(__ENV.DATA_PROFILE);
const products = buildProducts(profile);

export const requestFailureRate = new Rate('request_failure_rate');
export const operationDuration = new Trend('operation_duration', true);
export const createdProductsCounter = new Counter('created_products');
export const updatedProductsCounter = new Counter('updated_products');
export const deletedProductsCounter = new Counter('deleted_products');

export const options = {
  vus: profile.vus,
  iterations: profile.iterations,
  thresholds: {
    checks: ['rate>0.99'],
    request_failure_rate: ['rate<0.01'],
    operation_duration: [`p(95)<${profile.durationP95Ms}`],
    created_products: [`count>=${profile.productCount}`],
    updated_products: [`count>=${profile.updateSampleSize}`],
    deleted_products: ['count>=1'],
  },
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)'],
};

const grpcClient = target.protocol === 'grpc' ? new grpc.Client() : null;
if (grpcClient) {
  grpcClient.load(['../shared/proto'], 'product/v1/product.proto');
}

function recordResult(operationName, ok, durationMs) {
  requestFailureRate.add(ok ? 0 : 1, { operation: operationName, target: target.id });
  operationDuration.add(durationMs, { operation: operationName, target: target.id });
}

function recordCheck(operationName, result, checksMap) {
  const ok = check(result, checksMap, { operation: operationName, target: target.id });
  return ok;
}

function ensureValue(operationName, value, message) {
  if (value === undefined || value === null) {
    throw new Error(`${operationName}: ${message}`);
  }

  return value;
}

function getStockQuantity(product) {
  return product.stockQuantity ?? product.stock_quantity;
}

function getImages(product) {
  return Array.isArray(product.images) ? product.images : [];
}

function arraysEqual(left, right) {
  if (!Array.isArray(left) || !Array.isArray(right) || left.length !== right.length) {
    return false;
  }

  return left.every((value, index) => value === right[index]);
}

function productsMatch(actual, expected) {
  if (!actual || !expected) {
    return false;
  }

  return (
    actual.name === expected.name &&
    actual.description === expected.description &&
    actual.category === expected.category &&
    Number(actual.price) === Number(expected.price) &&
    Number(getStockQuantity(actual)) === Number(expected.stockQuantity) &&
    arraysEqual(getImages(actual), expected.images)
  );
}

function restHeaders() {
  return {
    headers: {
      'Content-Type': 'application/json',
    },
  };
}

function httpDeleteAllowingNotFound(path) {
  const startedAt = Date.now();
  const response = http.del(`${target.baseUrl}${path}`);
  const ok = recordCheck('cleanup_delete_all', response, {
    'cleanup delete-all status ok': (res) => res.status === 204 || res.status === 404,
  });
  recordResult('cleanup_delete_all', ok, Date.now() - startedAt);
}

function graphQLRequest(operationName, query, variables = {}) {
  const startedAt = Date.now();
  const response = http.post(
    target.baseUrl,
    JSON.stringify({ query, variables }),
    restHeaders(),
  );
  const ok = recordCheck(operationName, response, {
    'graphql status is 200': (res) => res.status === 200,
    'graphql has no errors': (res) => {
      const body = JSON.parse(res.body);
      return !body.errors;
    },
  });
  recordResult(operationName, ok, Date.now() - startedAt);
  return JSON.parse(response.body);
}

function grpcInvoke(operationName, methodName, payload) {
  const startedAt = Date.now();
  const response = grpcClient.invoke(methodName, payload);
  const ok = recordCheck(operationName, response, {
    'grpc status is OK': (res) => res && res.status === grpc.StatusOK,
  });
  recordResult(operationName, ok, Date.now() - startedAt);
  return response;
}

function ensureGrpcConnection() {
  if (grpcClient) {
    grpcClient.connect(target.address, { plaintext: true });
  }
}

function cleanupBeforeRun() {
  if (target.protocol === 'rest') {
    httpDeleteAllowingNotFound('/Product/all');
    return;
  }

  if (target.protocol === 'graphql') {
    graphQLRequest('cleanup_delete_all', 'mutation { deleteAllProducts }');
    return;
  }

  ensureGrpcConnection();
  grpcInvoke('cleanup_delete_all', 'product.v1.ProductService/DeleteAllProducts', {});
}

function cleanupAfterRun() {
  if (target.protocol === 'rest') {
    const startedAt = Date.now();
    const response = http.del(`${target.baseUrl}/Product/all`);
    const ok = recordCheck('final_delete_all', response, {
      'final delete-all status ok': (res) => res.status === 204 || res.status === 404,
    });
    recordResult('final_delete_all', ok, Date.now() - startedAt);
    if (ok && response.status === 204) {
      deletedProductsCounter.add(1, { target: target.id });
    }
    return;
  }

  if (target.protocol === 'graphql') {
    const body = graphQLRequest('final_delete_all', 'mutation { deleteAllProducts }');
    if (body.data?.deleteAllProducts) {
      deletedProductsCounter.add(1, { target: target.id });
    }
    return;
  }

  ensureGrpcConnection();
  const response = grpcInvoke('final_delete_all', 'product.v1.ProductService/DeleteAllProducts', {});
  if (response.message?.deleted) {
    deletedProductsCounter.add(1, { target: target.id });
  }
}

function createProductRest(product) {
  const startedAt = Date.now();
  const response = http.post(`${target.baseUrl}/Product`, JSON.stringify(product), restHeaders());
  const ok = recordCheck('create_product', response, {
    'rest create status is 200': (res) => res.status === 200,
    'rest create echoes payload': (res) => productsMatch(JSON.parse(res.body), product),
  });
  recordResult('create_product', ok, Date.now() - startedAt);
  const body = JSON.parse(response.body);
  createdProductsCounter.add(1, { target: target.id });
  return ensureValue('create_product', body.id, 'resposta REST sem id');
}

function getAllRest(expectedProducts) {
  const startedAt = Date.now();
  const response = http.get(`${target.baseUrl}/Product/all`);
  const ok = recordCheck('get_all_products', response, {
    'rest get-all status is 200': (res) => res.status === 200,
    'rest get-all expected count': (res) => JSON.parse(res.body).length === expectedProducts.length,
    'rest get-all includes full payloads': (res) => {
      const body = JSON.parse(res.body);
      return expectedProducts.every((expected) => body.some((item) => productsMatch(item, expected)));
    },
  });
  recordResult('get_all_products', ok, Date.now() - startedAt);
}

function getByIdRest(id, expectedProduct) {
  const startedAt = Date.now();
  const response = http.get(`${target.baseUrl}/Product?id=${id}`);
  const ok = recordCheck('get_product_by_id', response, {
    'rest get-by-id status is 200': (res) => res.status === 200,
    'rest get-by-id expected payload': (res) => productsMatch(JSON.parse(res.body), expectedProduct),
  });
  recordResult('get_product_by_id', ok, Date.now() - startedAt);
}

function updateProductRest(id, product) {
  const startedAt = Date.now();
  const response = http.put(
    `${target.baseUrl}/Product/${id}`,
    JSON.stringify(product),
    restHeaders(),
  );
  const ok = recordCheck('update_product', response, {
    'rest update status is 200': (res) => res.status === 200,
    'rest update expected payload': (res) => productsMatch(JSON.parse(res.body), product),
  });
  recordResult('update_product', ok, Date.now() - startedAt);
  updatedProductsCounter.add(1, { target: target.id });
}

function createProductGraphQL(product) {
  const body = graphQLRequest(
    'create_product',
    `
      mutation ($input: CreateProductInput!) {
        createProduct(input: $input) {
          id
          name
          description
          category
          images
          price
          stockQuantity
        }
      }
    `,
    { input: product },
  );
  const ok = check(body, {
    'graphql create echoes payload': (data) => productsMatch(data.data?.createProduct, product),
  }, { operation: 'create_product_payload', target: target.id });
  requestFailureRate.add(ok ? 0 : 1, { operation: 'create_product_payload', target: target.id });
  createdProductsCounter.add(1, { target: target.id });
  return ensureValue('create_product', body.data?.createProduct?.id, 'resposta GraphQL sem id');
}

function getAllGraphQL(expectedProducts) {
  const body = graphQLRequest(
    'get_all_products',
    `
      query {
        allProducts {
          id
          name
          description
          category
          images
          price
          stockQuantity
        }
      }
    `,
  );

  const ok = check(body, {
    'graphql get-all expected count': (data) => data.data.allProducts.length === expectedProducts.length,
    'graphql get-all includes full payloads': (data) =>
      expectedProducts.every((expected) => data.data.allProducts.some((item) => productsMatch(item, expected))),
  }, { operation: 'get_all_products', target: target.id });
  requestFailureRate.add(ok ? 0 : 1, { operation: 'get_all_products_count', target: target.id });
}

function getByIdGraphQL(id, expectedProduct) {
  const body = graphQLRequest(
    'get_product_by_id',
    `
      query ($id: Int!) {
        productById(id: $id) {
          id
          name
          description
          category
          images
          price
          stockQuantity
        }
      }
    `,
    { id },
  );

  const ok = check(body, {
    'graphql get-by-id expected payload': (data) => productsMatch(data.data.productById, expectedProduct),
  }, { operation: 'get_product_by_id_name', target: target.id });
  requestFailureRate.add(ok ? 0 : 1, { operation: 'get_product_by_id_name', target: target.id });
}

function updateProductGraphQL(id, product) {
  const body = graphQLRequest(
    'update_product',
    `
      mutation ($id: Int!, $input: UpdateProductInput!) {
        updateProduct(id: $id, input: $input) {
          id
          name
          description
          category
          images
          price
          stockQuantity
        }
      }
    `,
    {
      id,
      input: {
        ...product,
        stockQuantity: product.stockQuantity,
      },
    },
  );

  const ok = check(body, {
    'graphql update expected payload': (data) => productsMatch(data.data.updateProduct, product),
  }, { operation: 'update_product_name', target: target.id });
  requestFailureRate.add(ok ? 0 : 1, { operation: 'update_product_name', target: target.id });
  updatedProductsCounter.add(1, { target: target.id });
}

function createProductGrpc(product) {
  const response = grpcInvoke('create_product', 'product.v1.ProductService/CreateProduct', {
    name: product.name,
    description: product.description,
    category: product.category,
    images: product.images,
    price: product.price,
    stock_quantity: product.stockQuantity,
  });
  const ok = check(response, {
    'grpc create echoes payload': (res) => productsMatch(res.message, product),
  }, { operation: 'create_product_payload', target: target.id });
  requestFailureRate.add(ok ? 0 : 1, { operation: 'create_product_payload', target: target.id });
  createdProductsCounter.add(1, { target: target.id });
  return ensureValue('create_product', response.message?.id, 'resposta gRPC sem id');
}

function getAllGrpc(expectedProducts) {
  const response = grpcInvoke('get_all_products', 'product.v1.ProductService/GetAllProducts', {});
  const ok = check(response, {
    'grpc get-all expected count': (res) => res.message.products.length === expectedProducts.length,
    'grpc get-all includes full payloads': (res) =>
      expectedProducts.every((expected) => res.message.products.some((item) => productsMatch(item, expected))),
  }, { operation: 'get_all_products_count', target: target.id });
  requestFailureRate.add(ok ? 0 : 1, { operation: 'get_all_products_count', target: target.id });
}

function getByIdGrpc(id, expectedProduct) {
  const response = grpcInvoke('get_product_by_id', 'product.v1.ProductService/GetProductById', { id });
  const ok = check(response, {
    'grpc get-by-id expected payload': (res) => productsMatch(res.message, expectedProduct),
  }, { operation: 'get_product_by_id_name', target: target.id });
  requestFailureRate.add(ok ? 0 : 1, { operation: 'get_product_by_id_name', target: target.id });
}

function updateProductGrpc(id, product) {
  const response = grpcInvoke('update_product', 'product.v1.ProductService/UpdateProduct', {
    id,
    name: product.name,
    description: product.description,
    category: product.category,
    images: product.images,
    price: product.price,
    stock_quantity: product.stockQuantity,
  });
  const ok = check(response, {
    'grpc update expected payload': (res) => productsMatch(res.message, product),
  }, { operation: 'update_product_name', target: target.id });
  requestFailureRate.add(ok ? 0 : 1, { operation: 'update_product_name', target: target.id });
  updatedProductsCounter.add(1, { target: target.id });
}

function runRestScenario() {
  const ids = products.map((product) => createProductRest(product));
  getAllRest(products);

  const readIndexes = pickIndexes(ids.length, profile.readSampleSize);
  for (const index of readIndexes) {
    getByIdRest(ids[index], products[index]);
  }

  const updateIndexes = pickIndexes(ids.length, profile.updateSampleSize);
  for (const index of updateIndexes) {
    const updated = buildUpdatedProduct(products[index], index);
    updateProductRest(ids[index], updated);
    products[index] = updated;
  }
}

function runGraphQLScenario() {
  const ids = products.map((product) => createProductGraphQL(product));
  getAllGraphQL(products);

  const readIndexes = pickIndexes(ids.length, profile.readSampleSize);
  for (const index of readIndexes) {
    getByIdGraphQL(ids[index], products[index]);
  }

  const updateIndexes = pickIndexes(ids.length, profile.updateSampleSize);
  for (const index of updateIndexes) {
    const updated = buildUpdatedProduct(products[index], index);
    updateProductGraphQL(ids[index], updated);
    products[index] = updated;
  }
}

function runGrpcScenario() {
  ensureGrpcConnection();
  const ids = products.map((product) => createProductGrpc(product));
  getAllGrpc(products);

  const readIndexes = pickIndexes(ids.length, profile.readSampleSize);
  for (const index of readIndexes) {
    getByIdGrpc(ids[index], products[index]);
  }

  const updateIndexes = pickIndexes(ids.length, profile.updateSampleSize);
  for (const index of updateIndexes) {
    const updated = buildUpdatedProduct(products[index], index);
    updateProductGrpc(ids[index], updated);
    products[index] = updated;
  }
}

export function setup() {
  cleanupBeforeRun();
  return { targetId: target.id, profile: profile.label };
}

export default function () {
  if (target.protocol === 'rest') {
    runRestScenario();
    return;
  }

  if (target.protocol === 'graphql') {
    runGraphQLScenario();
    return;
  }

  runGrpcScenario();
}

export function teardown() {
  cleanupAfterRun();
  if (grpcClient) {
    grpcClient.close();
  }
}

function escapeHtml(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;');
}

function formatMetricValue(value) {
  if (typeof value !== 'number' || Number.isNaN(value)) {
    return '-';
  }

  if (Number.isInteger(value)) {
    return String(value);
  }

  return value.toFixed(4);
}

function buildMetricRows(data) {
  const metricNames = [
    'checks',
    'request_failure_rate',
    'operation_duration',
    'created_products',
    'updated_products',
    'deleted_products',
    'iterations',
    'iteration_duration',
  ];

  return metricNames
    .map((name) => {
      const metric = data.metrics[name];
      if (!metric) {
        return '';
      }

      const values = metric.values || {};
      const details = Object.entries(values)
        .map(([key, value]) => `${escapeHtml(key)}: ${escapeHtml(formatMetricValue(value))}`)
        .join('<br>');

      return `
        <tr>
          <td>${escapeHtml(name)}</td>
          <td>${escapeHtml(metric.type || '-')}</td>
          <td>${details || '-'}</td>
          <td>${escapeHtml((metric.thresholds && Object.keys(metric.thresholds).join(', ')) || '-')}</td>
        </tr>
      `;
    })
    .join('');
}

function buildHtmlReport(data, failureRate) {
  const title = `k6 report - ${target.id} - ${profile.label}`;
  const createdAt = new Date().toISOString();

  return `<!doctype html>
<html lang="pt-BR">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>${escapeHtml(title)}</title>
    <style>
      :root {
        color-scheme: light;
        --bg: #f4efe6;
        --panel: #fffdf8;
        --ink: #1f2937;
        --muted: #6b7280;
        --line: #d6cfc2;
        --accent: #b45309;
        --ok: #166534;
        --warn: #b91c1c;
      }
      * { box-sizing: border-box; }
      body {
        margin: 0;
        font-family: "IBM Plex Sans", "Segoe UI", sans-serif;
        background: linear-gradient(180deg, #f7f2ea 0%, var(--bg) 100%);
        color: var(--ink);
      }
      main {
        max-width: 1100px;
        margin: 0 auto;
        padding: 32px 20px 48px;
      }
      .hero, .panel {
        background: var(--panel);
        border: 1px solid var(--line);
        border-radius: 18px;
        box-shadow: 0 12px 30px rgba(64, 43, 17, 0.08);
      }
      .hero {
        padding: 24px;
        margin-bottom: 20px;
      }
      h1, h2 { margin: 0 0 12px; }
      h1 { font-size: 2rem; }
      h2 { font-size: 1.1rem; }
      .muted { color: var(--muted); }
      .grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
        gap: 12px;
        margin-top: 20px;
      }
      .card {
        padding: 16px;
        border-radius: 14px;
        border: 1px solid var(--line);
        background: #fffaf1;
      }
      .label {
        display: block;
        color: var(--muted);
        font-size: 0.85rem;
        margin-bottom: 8px;
      }
      .value {
        font-size: 1.35rem;
        font-weight: 700;
      }
      .ok { color: var(--ok); }
      .warn { color: var(--warn); }
      .panel {
        padding: 20px;
      }
      table {
        width: 100%;
        border-collapse: collapse;
      }
      th, td {
        padding: 12px 10px;
        text-align: left;
        vertical-align: top;
        border-top: 1px solid var(--line);
      }
      th {
        border-top: none;
        color: var(--muted);
        font-weight: 600;
      }
      @media (max-width: 720px) {
        h1 { font-size: 1.5rem; }
        th, td { font-size: 0.92rem; }
      }
    </style>
  </head>
  <body>
    <main>
      <section class="hero">
        <h1>${escapeHtml(title)}</h1>
        <p class="muted">Gerado em ${escapeHtml(createdAt)}</p>
        <div class="grid">
          <div class="card">
            <span class="label">Target</span>
            <span class="value">${escapeHtml(target.id)}</span>
          </div>
          <div class="card">
            <span class="label">Perfil</span>
            <span class="value">${escapeHtml(profile.label)}</span>
          </div>
          <div class="card">
            <span class="label">Falha</span>
            <span class="value ${failureRate > 0.01 ? 'warn' : 'ok'}">${escapeHtml(formatMetricValue(failureRate * 100))}%</span>
          </div>
          <div class="card">
            <span class="label">p95 operacao</span>
            <span class="value">${escapeHtml(formatMetricValue(data.metrics.operation_duration?.values?.['p(95)']))} ms</span>
          </div>
        </div>
      </section>
      <section class="panel">
        <h2>Metricas principais</h2>
        <table>
          <thead>
            <tr>
              <th>Metrica</th>
              <th>Tipo</th>
              <th>Valores</th>
              <th>Thresholds</th>
            </tr>
          </thead>
          <tbody>
            ${buildMetricRows(data)}
          </tbody>
        </table>
      </section>
    </main>
  </body>
</html>`;
}

export function handleSummary(data) {
  const failureRate = data.metrics.request_failure_rate?.values?.rate ?? 0;
  const status = failureRate > 0.01 ? 'com falhas' : 'sem falhas';
  const htmlReportFile = __ENV.HTML_REPORT_FILE;
  const output = {
    stdout: `Suite finalizada para ${target.id} com perfil ${profile.label} (${status})\n`,
  };

  if (htmlReportFile) {
    output[htmlReportFile] = buildHtmlReport(data, failureRate);
  }

  return output;
}
