FROM golang:1.23.0-alpine3.20.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/api ./cmd/api
RUN go build -o /app/cli ./cmd/cli

FROM alpine:3.20.2

COPY --from=builder /app/api /app/api
COPY --from=builder /app/cli /app/cli

ENTRYPOINT ["/app/api"]
