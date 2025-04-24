package handlers

import (
	"api-gateway/config"
	"api-gateway/gateway/grpc/clients"
	"github.com/gin-gonic/gin"
)

func SetupUser(group *gin.RouterGroup, grpcClient *clients.Client, cfg *config.Config) {

	group.POST("/register", func(c *gin.Context) {
		RegisterUser(c, grpcClient)
	})

	group.POST("/login", func(c *gin.Context) {
		LoginUser(c, grpcClient)
	})

	group.GET("/profile", func(c *gin.Context) {
		GetProfile(c, grpcClient)
	})
}

func RegisterUser(c *gin.Context, grpcClient *clients.Client) {

}

func LoginUser(c *gin.Context, grpcClient *clients.Client) {

}

func GetProfile(c *gin.Context, grpcClient *clients.Client) {

}
