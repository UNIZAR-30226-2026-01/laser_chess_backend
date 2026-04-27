# Etapa 1: Compilación (usamos alpine solo para tener las herramientas de Go)
FROM golang:1.25.7-alpine AS builder
WORKDIR /app

# Instalamos certificados y SQLC (para replicar tu comando 'make sqlc')
RUN apk --no-cache add ca-certificates
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Copiamos dependencias primero
COPY go.mod go.sum ./
RUN go mod download

# Copiamos todo el código fuente
COPY . .

# Generamos el código de la BD (equivalente a 'make sqlc')
RUN sqlc generate

# Compilamos el backend apuntando a la ruta correcta (equivalente a 'make build')
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Etapa 2: Imagen ultra-ligera
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /main

EXPOSE 8080
CMD ["/main"]