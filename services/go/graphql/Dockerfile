FROM golang:1.22-alpine AS build
WORKDIR /src

COPY services/go/graphql/go.mod services/go/graphql/go.sum ./
RUN go mod download

COPY services/go/graphql/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/graphql-api ./src

FROM alpine:3.21 AS runtime
WORKDIR /app

RUN apk add --no-cache ca-certificates

ENV HTTP_ADDRESS=:8090
COPY --from=build /app/graphql-api .

ENTRYPOINT ["./graphql-api"]
