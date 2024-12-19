FROM golang:1.23.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/api ./cmd/api
RUN CGO_ENABLED=0 go build -o /app/cli ./cmd/cli
RUN CGO_ENABLED=0 go build -o /app/health ./cmd/health

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/api /app/api
COPY --from=builder /app/cli /app/cli
COPY --from=builder /app/health /app/health

ENTRYPOINT ["/app/api"]
