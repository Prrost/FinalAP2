package handlers

import (
	"net/http"
	"strconv"
	"time"

	"api-gateway/config"
	"api-gateway/gateway/Response"
	"api-gateway/gateway/grpc/clients"

	bookpb "github.com/Prrost/protoFinalAP2/books"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// SetupRouter — расширили, добавили /api/books
func SetupRouter(cfg *config.Config, grpcClient *clients.Client) *gin.Engine {
	router := gin.Default()

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
	SetupBooks(booksGroup, grpcClient)

	return router
}

// SetupBooks — новый внутри того же файла
func SetupBooks(group *gin.RouterGroup, grpcClient *clients.Client) {
	// GET /api/books
	group.GET("", func(c *gin.Context) {
		res, err := grpcClient.BookClient.ListBooks(c.Request.Context(), &emptypb.Empty{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response.Err{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, res.Books)
	})

	// GET /api/books/:id
	group.GET("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response.Err{Error: "invalid book id"})
			return
		}
		book, err := grpcClient.BookClient.GetBook(c.Request.Context(), &bookpb.BookId{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response.Err{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, book)
	})

	// POST /api/books
	group.POST("", func(c *gin.Context) {
		var in bookpb.Book
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, Response.Err{Error: err.Error()})
			return
		}
		created, err := grpcClient.BookClient.CreateBook(c.Request.Context(), &in)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response.Err{Error: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, created)
	})

	// PUT /api/books/:id
	group.PUT("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response.Err{Error: "invalid book id"})
			return
		}
		var in bookpb.Book
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, Response.Err{Error: err.Error()})
			return
		}
		in.Id = id
		updated, err := grpcClient.BookClient.UpdateBook(c.Request.Context(), &in)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response.Err{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)
	})

	// DELETE /api/books/:id
	group.DELETE("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response.Err{Error: "invalid book id"})
			return
		}
		_, err = grpcClient.BookClient.DeleteBook(c.Request.Context(), &bookpb.BookId{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response.Err{Error: err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
