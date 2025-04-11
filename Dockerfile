# Etapa de build
FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
# Aqui está o segredo do sucesso:
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tinyliar ./cmd/main.go

# Etapa final
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/tinyliar .

EXPOSE 8888
CMD ["./tinyliar"]
