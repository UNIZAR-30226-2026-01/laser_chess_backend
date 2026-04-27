# Etapa 1: Compilación (usamos alpine solo para tener las herramientas de Go)
FROM golang:1.25.7-alpine AS builder
WORKDIR /app

# Instalamos los certificados por si tu backend hace peticiones HTTPS hacia afuera
RUN apk --no-cache add ca-certificates

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Compilamos el binario. 
# CGO_ENABLED=0 es OBLIGATORIO para 'scratch' para que sea un binario 100% estático.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Etapa 2: El vacío absoluto
FROM scratch

# Copiamos los certificados SSL de la etapa anterior (vital si llamas a APIs externas)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiamos nuestro binario ya compilado
COPY --from=builder /app/main /main

# Exponemos el puerto
EXPOSE 8080

# Ejecutamos el binario directamente
CMD ["/main"]