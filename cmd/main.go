package main

import (
	"context"
	"log"
	"os"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/rt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	// Cargar las variables de .env
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatalln("Error: DB_URL no encontrada")
	}

	// == INICIALIZAR BDD ===============================================

	// Crear pool de conexiones con la db
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalln("Error: No se pudo conectar con la base de datos")
	}
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalln("Error: No hay conexión real con la DB:", err)
	}
	defer dbPool.Close()

	// Inicializar store de sqlc
	store := db.NewStore(dbPool)

	// == INICIALIZAR HUBS ===============================================

	// Inicializar el registro de partidas activas
	matchRegistry := rt.NewMatchRegistry()

	// Inicializar el hub privado
	privateHub := rt.NewPrivateHub(matchRegistry)

	// Inicializar el hub publico
	publicHub := rt.NewPublicHub(matchRegistry)

	// == INICIALIZAR ROUTER =============================================

	// Inicializar router de gin
	router := api.SetupRouter(store, matchRegistry, privateHub, publicHub)

	// Ejecutar el router en el puerto que se le diga (8080 por defecto)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
