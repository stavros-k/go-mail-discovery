FROM golang:1.23.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/api ./cmd/api
RUN CGO_ENABLED=0 go build -o /app/cli ./cmd/cli

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/api /app/api
COPY --from=builder /app/cli /app/cli

ENTRYPOINT ["/app/api"]
