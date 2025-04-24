package handlers

import (
	"api-gateway/config"
	"api-gateway/gateway/grpc/clients"
	"api-gateway/gateway/middleware"
	"github.com/gin-gonic/gin"
)

// just an example of real thing
func SetupRouter(cfg *config.Config, grpcClient *clients.Client) *gin.Engine {
	router := gin.Default()

	apiGroup := router.Group("/api")

	productsGroup := apiGroup.Group("/products")
	ordersGroup := apiGroup.Group("/orders")
	userGroup := apiGroup.Group("/user")

	productsGroup.Use(middleware.AuthMiddleware(cfg))
	ordersGroup.Use(middleware.AuthMiddleware(cfg))

	SetupUser(userGroup, grpcClient, cfg)

	return router
}
