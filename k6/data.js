export function buildProducts(profile) {
  return Array.from({ length: profile.productCount }, (_, index) => {
    const number = index + 1;
    return {
      name: `K6 Product ${String(number).padStart(4, '0')}`,
      description: `Dataset ${profile.label} item ${number}`,
      price: Number((number * 3.75 + 10).toFixed(2)),
      stockQuantity: number * 2,
    };
  });
}

export function pickIndexes(total, desired) {
  if (total === 0 || desired === 0) {
    return [];
  }

  if (desired >= total) {
    return Array.from({ length: total }, (_, index) => index);
  }

  const step = total / desired;
  const indexes = [];

  for (let i = 0; i < desired; i += 1) {
    const value = Math.min(total - 1, Math.floor(i * step));
    if (!indexes.includes(value)) {
      indexes.push(value);
    }
  }

  return indexes;
}

export function buildUpdatedProduct(product, index) {
  return {
    name: `${product.name} Updated`,
    description: `${product.description} Updated ${index + 1}`,
    price: Number((product.price + 1.25).toFixed(2)),
    stockQuantity: product.stockQuantity + 5,
  };
}
