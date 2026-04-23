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

## Carga com k6

O teste de carga fica em [k6/api-load.js](/home/felipe/Documentos/Programming/TCC/TrabalhoPratico/k6/api-load.js) e usa o mesmo dataset deterministico para todas as APIs.

Perfis suportados:

- `single`: 1 produto
- `small`: 25 produtos
- `medium`: 100 produtos
- `large`: 500 produtos

Aliases em portugues aceitos pelo runner:

- `unico`
- `pequeno`
- `medio`
- `grande`

Runner sequencial por projeto:

```bash
./scripts/run-k6-suite.sh single
./scripts/run-k6-suite.sh small
./scripts/run-k6-suite.sh medium
./scripts/run-k6-suite.sh large
```

Tambem e possivel executar um unico projeto:

```bash
./scripts/run-k6-suite.sh medium go-graphql
```

Comportamento do runner:

- executa um projeto por vez
- chama `delete all` no inicio da rodada
- executa `create`, `get by id`, `get all` e `update`
- chama `delete all` ao final antes de seguir para o proximo projeto
- grava o resumo JSON e o relatorio HTML de cada execucao em `result/k6/`

Metricas validadas:

- `checks > 99%`
- `request_failure_rate < 1%`
- `operation_duration` com `p(95)` por perfil
- contagem minima de criacoes e atualizacoes esperadas

## Observacoes

- O banco eh inicializado com `database/init.sql` e `database/seed.sql`.
- O script de validacao usa o binario local em `.tools/bin/grpcurl`.
- O servico `csharp-grpc` usa a variavel `GRPC_PORT` para bind da porta dentro do container.
- Cada API sobe em um container independente, sem compartilhar processo, porta ou runtime com as demais.
