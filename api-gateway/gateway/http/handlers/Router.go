package handlers

import (
	"time"

	"api-gateway/config"
	"api-gateway/gateway/grpc/clients"
	"api-gateway/gateway/middleware" // путь к твоему middleware
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, grpcClient *clients.Client) *gin.Engine {
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Front},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           100 * time.Hour,
	}))

	api := router.Group("/api")

	// User
	userGroup := api.Group("/user")
	SetupUser(userGroup, grpcClient, cfg)

	// Books
	booksGroup := api.Group("/books")
	// вот здесь — правильно подключаем твой AuthMiddleware
	booksGroup.Use(middleware.AuthMiddleware(cfg))
	// и сразу регистрируем CRUD без повторного Use
	SetupBooks(booksGroup, grpcClient)

	return router
}
