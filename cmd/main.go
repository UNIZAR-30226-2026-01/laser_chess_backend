package main

import (
	"context"
	"log"
	"os"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/item"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/login"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/match"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/gin-gonic/gin"
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

	// Inicializar router de gin
	router := gin.Default()

	// Inicializar store de sqlc
	store := db.NewStore(dbPool)

	// Crear handlers y services
	loginService := login.NewService(store)
	loginHandler := login.NewHandler(loginService)

	accountService := account.NewService(store)
	accountHandler := account.NewHandler(accountService)

	matchService := match.NewService(store)
	matchHandler := match.NewHandler(matchService)

	itemService := item.NewService(store)
	itemHandler := item.NewHandler(itemService)

	// Establecer las rutas de las peticiones http por grupos

	// Login routes

	router.POST("/login", loginHandler.Login)

	// Account routes
	{
		accountRoute := router.Group("/account")
		accountRoute.POST("", accountHandler.Create)
		accountRoute.GET("/:id", accountHandler.GetByID)
		accountRoute.POST("/update", accountHandler.Update)
		accountRoute.DELETE("/delete/:id", accountHandler.Delete)
	}

	// Match routes
	{
		matchRoute := router.Group("/match")
		matchRoute.POST("", matchHandler.CreateMatch)
		matchRoute.GET("/:matchID", matchHandler.GetMatch)
		matchRoute.GET("/history/:userID", matchHandler.GetUserHistory)
	}

	// Item routes
	{
		itemRoute := router.Group("/item")
		itemRoute.POST("", itemHandler.CreateItemOwner)
		itemRoute.GET("/inventory/:userID", itemHandler.GetUserItems)
		itemRoute.GET("/:itemID", itemHandler.GetShopItem)
	}

	// Ejecutar el router en el puerto que se le diga (8080 por defecto)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
