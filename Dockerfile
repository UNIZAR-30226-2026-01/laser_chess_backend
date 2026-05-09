# Etapa 1: Compilación (usamos alpine solo para tener las herramientas de Go)
FROM golang:1.25.7-alpine AS builder
WORKDIR /app

ENV CGO_ENABLED=0

# Instalamos certificados y SQLC (para replicar tu comando 'make sqlc')
RUN apk --no-cache add ca-certificates wget tar make gcc musl-dev
RUN wget https://downloads.sqlc.dev/sqlc_1.25.0_linux_amd64.tar.gz && \
    tar -xvzf sqlc_1.25.0_linux_amd64.tar.gz && \
    mv sqlc /usr/local/bin/

# Copiamos dependencias primero
COPY go.mod go.sum ./
RUN go mod download

# Copiamos todo el código fuente
COPY . .

# Generamos el código de la BD (equivalente a 'make sqlc')
RUN sqlc generate

# Compilamos el backend apuntando a la ruta correcta (equivalente a 'make build')
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o seeder-bin ./cmd/seed/main.go

# Etapa 2: Imagen ultra-ligera
FROM scratch AS backend
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /main

EXPOSE 8080
CMD ["/main"]

FROM alpine:latest AS seeder
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/seeder-bin /seeder
# El seeder no suele necesitar EXPOSE porque solo escribe en la BD y termina
CMD ["/seeder"]

FROM postgres:16-alpine AS database
# Copiamos el esquema a la base de datos y hacemos una imagen con el esquema aplicado ya
COPY internal/db/schema.sql /docker-entrypoint-initdb.d/1-schema.sql