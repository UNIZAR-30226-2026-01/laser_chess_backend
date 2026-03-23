package api

// Fichero que se encarga de inicializar todas las rutas de la api

import (
	"os"
	"time"

	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/account"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/friendship"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/item"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/login"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/match"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/middleware"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/api/rating"
	db "github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/db/sqlc"
	"github.com/UNIZAR-30226-2026-01/laser_chess_backend/internal/rt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Hace el setup del router
// Crea handlers y services, y mapea endpoints
// Tambien activa el middleware del jwt
func SetupRouter(store *db.Store,
	registry *rt.MatchRegistry,
	privateHub *rt.PrivateHub,
) *gin.Engine {

	router := gin.Default()

	// Cross-Origin Resource Sharing para conexion con Angular
	router.Use(cors.New(cors.Config{
		// Habrá que cambiar esto por la url real en produccion
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Crear handlers y services
	loginService := login.NewService(store)
	loginHandler := login.NewHandler(loginService)

	accountService := account.NewService(store)
	accountHandler := account.NewHandler(accountService)

	matchService := match.NewService(store)
	matchHandler := match.NewHandler(matchService)

	itemService := item.NewService(store)
	itemHandler := item.NewHandler(itemService)

	ratingService := rating.NewService(store)
	ratingHandler := rating.NewHandler(ratingService)

	friendshipService := friendship.NewService(store)
	friendshipHandler := friendship.NewHandler(friendshipService, accountService)

	// Establecer las rutas de las peticiones http por grupos

	// == RUTAS PUBLICAS ==============================================

	router.POST("/login", loginHandler.Login)
	router.POST("/register", accountHandler.Create)
	router.POST("/refresh", loginHandler.Refresh)

	// == RUTAS PRIVADAS ===============================================

	protected := router.Group("/api")

	// Si la variable de entorno "PROTECTED" == "FALSE" no se usa el middleware
	// de seguridad y se pueden hacer pruebas sin preocuparse por JWTs

	if os.Getenv("PROTECTED") != "FALSE" {
		protected.Use(middleware.AuthMiddleware())
	}

	// Account routes
	{
		accountRoute := protected.Group("/account")
		accountRoute.GET("/:id", accountHandler.GetByID)
		accountRoute.POST("/update", accountHandler.Update)
		accountRoute.DELETE("/delete/", accountHandler.Delete)
	}

	// Match routes
	{
		matchRoute := protected.Group("/match")
		// Probablemente las partidas las acabe creando la app, a si
		// que creo que no se usara el POST (la funcion del service si)
		matchRoute.POST("", matchHandler.CreateMatch)
		matchRoute.GET("/:matchID", matchHandler.GetMatch)
		matchRoute.GET("/history/:userID", matchHandler.GetUserHistory)
	}

	// Item routes
	{
		itemRoute := protected.Group("/item")
		itemRoute.POST("", itemHandler.CreateItemOwner)
		itemRoute.GET("/inventory", itemHandler.GetUserItems)
		itemRoute.GET("/:itemID", itemHandler.GetShopItem)
		// Probablemente habrá que meter una de GET todos los items
	}

	// Rating routes
	{
		ratingRoute := protected.Group("/rating")
		ratingRoute.GET("/:userID", ratingHandler.GetAllElos)
		ratingRoute.GET("/:userID/blitz", ratingHandler.GetBlitzElo)
		ratingRoute.GET("/:userID/bullet", ratingHandler.GetBulletElo)
		ratingRoute.GET("/:userID/rapid", ratingHandler.GetRapidElo)
		ratingRoute.GET("/:userID/classic", ratingHandler.GetClassicElo)
	}

	// Friendship routes
	{
		friendshipRoute := protected.Group("/friendship")
		friendshipRoute.POST("", friendshipHandler.Create)
		friendshipRoute.GET("", friendshipHandler.GetUserFriendships)

		friendshipRoute.GET("/sent",
			friendshipHandler.GetUserPendingSentFriendships)

		friendshipRoute.GET("/pending",
			friendshipHandler.GetUserPendingReceivedFriendships)

		friendshipRoute.GET("/:user2Username", friendshipHandler.GetFriendshipStatus)
		friendshipRoute.PUT("/:user2Username", friendshipHandler.AcceptFriendship)
		friendshipRoute.DELETE("/:user2Username", friendshipHandler.DeleteFriendship)
	}

	// Endpoints de websockets
	privateHandler := rt.NewPrivateHandler(privateHub, registry, accountService)

	{
		rtRoute := protected.Group("/rt/")
		rtRoute.GET("/challenge", privateHandler.Challenge)
		rtRoute.GET("/challenge/accept", privateHandler.AcceptChallenge)
		rtRoute.GET("/challenges", privateHandler.GetChallenges)
	}

	return router
}
