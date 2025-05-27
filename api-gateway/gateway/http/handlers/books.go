package handlers

import (
	"net/http"
	"strconv"

	"api-gateway/gateway/Response"
	"api-gateway/gateway/grpc/clients"
	bookpb "github.com/Prrost/protoFinalAP2/books"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func SetupBooks(group *gin.RouterGroup, grpcClient *clients.Client) {
	handleErr := func(c *gin.Context, err error) {
		if err == nil {
			return
		}
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, Response.Err{Error: st.Message()})
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, Response.Err{Error: st.Message()})
		case codes.AlreadyExists:
			c.JSON(http.StatusConflict, Response.Err{Error: st.Message()})
		default:
			c.JSON(http.StatusInternalServerError, Response.Err{Error: st.Message()})
		}
	}

	// List
	group.GET("", func(c *gin.Context) {
		res, err := grpcClient.BookClient.ListBooks(c.Request.Context(), &emptypb.Empty{})
		if err != nil {
			handleErr(c, err)
			return
		}
		c.JSON(http.StatusOK, res.Books)
	})

	// Get by ID
	group.GET("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response.Err{Error: "invalid book id"})
			return
		}
		book, err := grpcClient.BookClient.GetBook(c.Request.Context(), &bookpb.BookId{Id: id})
		if err != nil {
			handleErr(c, err)
			return
		}
		c.JSON(http.StatusOK, book)
	})

	// Create
	group.POST("", func(c *gin.Context) {
		var in bookpb.Book
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, Response.Err{Error: err.Error()})
			return
		}
		created, err := grpcClient.BookClient.CreateBook(c.Request.Context(), &in)
		if err != nil {
			handleErr(c, err)
			return
		}
		c.JSON(http.StatusCreated, created)
	})

	// Update
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
			handleErr(c, err)
			return
		}
		c.JSON(http.StatusOK, updated)
	})

	// Delete
	group.DELETE("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response.Err{Error: "invalid book id"})
			return
		}
		_, err = grpcClient.BookClient.DeleteBook(c.Request.Context(), &bookpb.BookId{Id: id})
		if err != nil {
			handleErr(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	})
}
