package handlers

import (
	"api-gateway/config"
	"api-gateway/domain"
	"api-gateway/gateway/Response"
	"api-gateway/gateway/grpc/clients"
	"fmt"
	userpb "github.com/Prrost/protoFinalAP2/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
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
	var userInput domain.User

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, Response.Err{Error: fmt.Sprintf("invalid credentials: %s", err)})
		return
	}

	res, err := grpcClient.UserClient.RegisterUser(c.Request.Context(), &userpb.RegisterRequest{
		Email:    userInput.Email,
		Password: userInput.Password,
		IsAdmin:  userInput.IsAdmin,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.IndentedJSON(http.StatusBadRequest, Response.Err{Error: st.Message()})
			case codes.Internal:
				c.IndentedJSON(http.StatusInternalServerError, Response.Err{Error: st.Message()})
			case codes.AlreadyExists:
				c.IndentedJSON(http.StatusConflict, Response.Err{Error: st.Message()})
			}
		} else {
			c.IndentedJSON(http.StatusInternalServerError, Response.Err{Error: err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"id":      res.GetId(),
		"email":   userInput.Email,
		"isAdmin": userInput.IsAdmin,
	})

}

func LoginUser(c *gin.Context, grpcClient *clients.Client) {
	var userInput domain.User

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, Response.Err{Error: fmt.Sprintf("invalid credentials: %s", err)})
	}

	res, err := grpcClient.UserClient.AuthenticateUser(c.Request.Context(), &userpb.AuthRequest{
		Email:    userInput.Email,
		Password: userInput.Password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.IndentedJSON(http.StatusBadRequest, Response.Err{Error: st.Message()})
			case codes.Internal:
				c.IndentedJSON(http.StatusInternalServerError, Response.Err{Error: st.Message()})
			case codes.NotFound:
				c.IndentedJSON(http.StatusNotFound, Response.Err{Error: st.Message()})
			case codes.PermissionDenied:
				c.IndentedJSON(http.StatusForbidden, Response.Err{Error: st.Message()})
			}
		} else {
			c.IndentedJSON(http.StatusInternalServerError, Response.Err{Error: err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"token": res.GetToken(),
	})
}

func GetProfile(c *gin.Context, grpcClient *clients.Client) {
	id := c.Query("id")
	email := c.Query("email")

	res, err := grpcClient.UserClient.GetUserInfo(c.Request.Context(), &userpb.UserInfoRequest{
		Id:    id,
		Email: email,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				c.IndentedJSON(http.StatusBadRequest, Response.Err{Error: st.Message()})
			case codes.Internal:
				c.IndentedJSON(http.StatusInternalServerError, Response.Err{Error: st.Message()})
			case codes.NotFound:
				c.IndentedJSON(http.StatusNotFound, Response.Err{Error: st.Message()})
			case codes.PermissionDenied:
				c.IndentedJSON(http.StatusForbidden, Response.Err{Error: st.Message()})
			}
		} else {
			c.IndentedJSON(http.StatusInternalServerError, Response.Err{Error: err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"id":      res.GetId(),
		"email":   res.GetEmail(),
		"isAdmin": res.GetIsAdmin(),
	})

}
