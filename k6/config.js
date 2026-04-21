export const dataProfiles = {
  single: {
    label: 'single',
    productCount: 1,
    readSampleSize: 1,
    updateSampleSize: 1,
    vus: 1,
    iterations: 1,
    durationP95Ms: 1500,
  },
  small: {
    label: 'small',
    productCount: 25,
    readSampleSize: 5,
    updateSampleSize: 5,
    vus: 1,
    iterations: 1,
    durationP95Ms: 2000,
  },
  medium: {
    label: 'medium',
    productCount: 100,
    readSampleSize: 10,
    updateSampleSize: 10,
    vus: 1,
    iterations: 1,
    durationP95Ms: 3000,
  },
  large: {
    label: 'large',
    productCount: 500,
    readSampleSize: 25,
    updateSampleSize: 25,
    vus: 1,
    iterations: 1,
    durationP95Ms: 5000,
  },
};

const targetCatalog = {
  'csharp-rest': {
    id: 'csharp-rest',
    protocol: 'rest',
    baseUrl: 'http://localhost:5107',
  },
  'go-rest': {
    id: 'go-rest',
    protocol: 'rest',
    baseUrl: 'http://localhost:8080',
  },
  'rust-rest': {
    id: 'rust-rest',
    protocol: 'rest',
    baseUrl: 'http://localhost:8083',
  },
  'csharp-graphql': {
    id: 'csharp-graphql',
    protocol: 'graphql',
    baseUrl: 'http://localhost:5207/graphql',
  },
  'go-graphql': {
    id: 'go-graphql',
    protocol: 'graphql',
    baseUrl: 'http://localhost:8090/graphql',
  },
  'rust-graphql': {
    id: 'rust-graphql',
    protocol: 'graphql',
    baseUrl: 'http://localhost:8091/graphql',
  },
  'csharp-grpc': {
    id: 'csharp-grpc',
    protocol: 'grpc',
    address: 'localhost:5307',
  },
  'go-grpc': {
    id: 'go-grpc',
    protocol: 'grpc',
    address: 'localhost:50051',
  },
  'rust-grpc': {
    id: 'rust-grpc',
    protocol: 'grpc',
    address: 'localhost:50052',
  },
};

export function resolveDataProfile(input) {
  const normalized = String(input || 'single').toLowerCase();
  const aliases = {
    unico: 'single',
    unico_volume: 'single',
    pequeno: 'small',
    medio: 'medium',
    grande: 'large',
  };
  const key = aliases[normalized] || normalized;
  const profile = dataProfiles[key];

  if (!profile) {
    throw new Error(`Perfil de dados invalido: ${input}`);
  }

  return profile;
}

export function resolveTarget(input) {
  const key = String(input || '').toLowerCase();
  const target = targetCatalog[key];

  if (!target) {
    throw new Error(`Target invalido: ${input}`);
  }

  return target;
}

export function allTargetIds() {
  return Object.keys(targetCatalog);
}
