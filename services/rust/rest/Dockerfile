FROM rust:1.89-bookworm AS build
WORKDIR /src

COPY services/rust/rest/Cargo.toml services/rust/rest/Cargo.lock ./
COPY services/rust/rest/src ./src
RUN cargo build --release

FROM debian:bookworm-slim AS runtime
WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENV HTTP_HOST=0.0.0.0
ENV HTTP_PORT=8083
COPY --from=build /src/target/release/rest .

ENTRYPOINT ["./rest"]
