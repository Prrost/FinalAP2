package handlers

import (
	"api-gateway/config"
	"api-gateway/gateway/grpc/clients"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, grpcClient *clients.Client) *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api")

	userGroup := apiGroup.Group("/user")

	//example of middleware
	//productsGroup.Use(middleware.AuthMiddleware(cfg))

	SetupUser(userGroup, grpcClient, cfg)

	return router
}
