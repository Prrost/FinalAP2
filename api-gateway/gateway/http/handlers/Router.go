package handlers

import (
	"api-gateway/config"
	"api-gateway/gateway/grpc/clients"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func SetupRouter(cfg *config.Config, grpcClient *clients.Client) *gin.Engine {
	router := gin.Default()

	//Фикс проблемы CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.Front},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           100 * time.Hour,
	}))

	apiGroup := router.Group("/api")

	userGroup := apiGroup.Group("/user")

	//example of middleware
	//productsGroup.Use(middleware.AuthMiddleware(cfg))

	SetupUser(userGroup, grpcClient, cfg)

	return router
}
