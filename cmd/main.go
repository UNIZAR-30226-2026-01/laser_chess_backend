package main

import (
	"context"
	"log"
	"os"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/item"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/match"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/placeholder"
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

	// Inicializar queries de sqlc
	queries := db.New(dbPool)

	// Crear handlers y services
	placeholderService := placeholder.NewService(queries)
	placeholderHandler := placeholder.NewHandler(placeholderService)

	accountService := account.NewService(queries)
	accountHandler := account.NewHandler(accountService)

	matchService := match.NewService(queries)
	matchHandler := match.NewHandler(matchService)

	itemService := item.NewService(queries)
	itemHandler := item.NewHandler(itemService)

	// Establecer las rutas de las peticiones http por grupos
	{
		placeholderRoute := router.Group("/placeholder")
		placeholderRoute.POST("", placeholderHandler.CreatePlaceholder)
		placeholderRoute.GET("/:id", placeholderHandler.GetPlaceholder)
	}

	// Account routes
	{
		accountRoute := router.Group("/account")
		accountRoute.POST("", accountHandler.CreateAccount)
		accountRoute.GET("/:id", accountHandler.GetAccountByID)
		accountRoute.POST("/update", accountHandler.UpdateAccount)
		accountRoute.DELETE("/delete/:id", accountHandler.DeleteAccount)
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

	router.GET("/ping", getHello)

	// Ejecutar el router en el puerto que se le diga (8080 por defecto)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

func getHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
