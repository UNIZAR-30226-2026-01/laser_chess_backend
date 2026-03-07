# Laser Chess - Backend [![CI](https://github.com/UNIZAR-30226-2026-01/laser_chess_backend/actions/workflows/ci.yaml/badge.svg)](https://github.com/UNIZAR-30226-2026-01/laser_chess_backend/actions/workflows/ci.yaml)

Backend del juego de mesa online Laser Chess.

## Estructura del proyecto
```text
.
├── cmd
│   └── main.go                     # Inicializacion del sistema
├── docs
├── internal
│   ├── api                         # Endpoints de la api REST
│   │   ├── account
│   │   ├── apierror                # Gestion de errores
│   │   ├── ...
│   │   ├── middleware              # Middleware de auth
│   │   └── routes.go               # Inicializacion de la api REST
│   ├── auth                    
│   │   ├── passwords.go            # Gestion de contraseñas seguras
│   │   └── token.go                # Utilidades para JWTs
│   ├── db              
│   │   ├── initial_inserts.sql     # Datos iniciales bdd
│   │   ├── query                   # Directorio de queries sql
│   │   ├── schema.sql              # Creacion de tablas 
│   │   └── sqlc                    # Codigo generado por sqlc
│   │       ├── ...
│   │       └── store.go            # Objeto tipo repositorio
│   ├── game                        # Gestion de partidas laser chess
│   │   └── ...
│   └── rt                          # Gestion de comunicacion a tiempo real
│       └── ...
├── Makefile                        # Makefile para puesta en marcha
└── ...

