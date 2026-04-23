FROM golang:1.22-alpine AS build
WORKDIR /src

COPY services/go/grpc/go.mod services/go/grpc/go.sum ./
RUN go mod download

COPY services/go/grpc/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/grpc-api ./src

FROM alpine:3.21 AS runtime
WORKDIR /app

RUN apk add --no-cache ca-certificates

ENV GRPC_ADDRESS=:50051
COPY --from=build /app/grpc-api .

ENTRYPOINT ["./grpc-api"]
