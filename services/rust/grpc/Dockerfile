FROM rust:1.89-bookworm AS build
WORKDIR /src

COPY services/rust/grpc/Cargo.toml services/rust/grpc/Cargo.lock ./
COPY services/rust/grpc/build.rs ./
COPY services/rust/grpc/src ./src
COPY shared/proto ../../../shared/proto
RUN cargo build --release

FROM debian:bookworm-slim AS runtime
WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENV GRPC_HOST=0.0.0.0
ENV GRPC_PORT=50052
COPY --from=build /src/target/release/grpc .

ENTRYPOINT ["./grpc"]
