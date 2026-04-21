## Trabalho Pratico

Monorepo com implementacoes de `ProductService` em tres linguagens (`C#`, `Go` e `Rust`) e tres estilos de API (`REST`, `GraphQL` e `gRPC`), todos usando PostgreSQL.

## Subindo com Docker Compose

Na raiz do projeto:

```bash
docker compose up --build -d
```

Servicos expostos:

- `csharp-rest`: `http://localhost:5107/Product/all`
- `csharp-graphql`: `http://localhost:5207/graphql`
- `csharp-grpc`: `localhost:5307`
- `go-rest`: `http://localhost:8080/Product/all`
- `go-graphql`: `http://localhost:8090/graphql`
- `go-grpc`: `localhost:50051`
- `rust-rest`: `http://localhost:8083/Product/all`
- `rust-graphql`: `http://localhost:8091/graphql`
- `rust-grpc`: `localhost:50052`

O PostgreSQL do Compose fica acessivel apenas na rede interna (`postgres:5432`) para evitar conflito com um banco local ja rodando na maquina.

Para derrubar o ambiente:

```bash
docker compose down
```

Para derrubar removendo o volume do banco:

```bash
docker compose down -v
```

## Validacao rapida

Depois de subir os containers, execute:

```bash
./scripts/test-containers.sh
```

O script espera os 9 servicos ficarem disponiveis e valida:

- leitura de produtos seedados
- comportamento de fallback para produto inexistente
- endpoints REST
- queries GraphQL
- chamadas gRPC com `grpcurl`

## Observacoes

- O banco eh inicializado com `database/init.sql` e `database/seed.sql`.
- O script de validacao usa o binario local em `.tools/bin/grpcurl`.
- O servico `csharp-grpc` usa a variavel `GRPC_PORT` para bind da porta dentro do container.
- Cada API sobe em um container independente, sem compartilhar processo, porta ou runtime com as demais.
