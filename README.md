# Laser Chess - Backend

Backend REST y en tiempo real para el juego de mesa online **Laser Chess**, desarrollado en Go. 

El sistema gestiona la autenticación de usuarios, la persistencia de datos con PostgreSQL, la lógica del juego y las partidas multijugador en tiempo real.

---

## Tecnologías y Arquitectura

* **Lenguaje:** Go (Golang)
* **API REST:** Rutas y middlewares personalizados para autenticación (JWT) y gestión de cuentas.
* **Base de Datos:** PostgreSQL con **sqlc** para la generación de código Go a partir de consultas SQL de tipo seguro (*type-safe*).
* **Tiempo Real:** Sistema de comunicación bidireccional (`internal/rt`) para el estado de las partidas y movimientos.
* **Seguridad:** Hashing seguro de contraseñas y validación de tokens JWT.

---

## Estructura del Proyecto

.
├── cmd/
│   └── main.go                  # Punto de entrada e inicialización del servidor
├── internal/
│   ├── api/                     # Endpoints REST, middlewares de Auth y rutas
│   ├── auth/                    # Lógica de JWT y hashing de contraseñas
│   ├── db/                      # Esquema SQL, queries, datos iniciales y código generado por sqlc
│   ├── game/                    # Lógica de juego de Laser Chess (reglas, tablero, turnos)
│   └── rt/                      # Gestión de comunicación en tiempo real (WebSockets)
└── Makefile                     # Automatización de tareas (compilación, migración, ejecución)

---

## Puesta en Marcha

### Requisitos previos
* Go 1.20+ (o versión correspondiente)
* PostgreSQL

### Ejecución
Puedes desplegar la base de datos y lanzar el servidor usando el `Makefile`:

```bash
# Iniciar el servidor
make run
