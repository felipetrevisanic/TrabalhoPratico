FROM rust:1.89-bookworm AS build
WORKDIR /src

COPY services/rust/graphql/Cargo.toml services/rust/graphql/Cargo.lock ./
COPY services/rust/graphql/src ./src
RUN cargo build --release

FROM debian:bookworm-slim AS runtime
WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENV HTTP_HOST=0.0.0.0
ENV HTTP_PORT=8091
COPY --from=build /src/target/release/graphql .

ENTRYPOINT ["./graphql"]
