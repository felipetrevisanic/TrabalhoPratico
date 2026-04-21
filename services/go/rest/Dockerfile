FROM golang:1.22-alpine AS build
WORKDIR /src

COPY services/go/rest/go.mod services/go/rest/go.sum ./
RUN go mod download

COPY services/go/rest/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/rest-api ./cmd/api

FROM alpine:3.21 AS runtime
WORKDIR /app

RUN apk add --no-cache ca-certificates

ENV HTTP_ADDRESS=:8080
COPY --from=build /app/rest-api .

ENTRYPOINT ["./rest-api"]
