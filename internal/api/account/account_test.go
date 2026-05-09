package account

import (
	"context"
	"log"
	"os"

	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var accountService AccountService

func init() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatalln("Error: DB_URL no encontrada")
	}

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalln("Error: No se pudo conectar con la base de datos")
	}
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalln("Error: No hay conexión real con la DB:", err)
	}
	defer dbPool.Close()

	store := db.NewStore(dbPool)

	accountService = AccountService{store: store}
}
