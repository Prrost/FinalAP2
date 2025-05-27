package httpserver

import (
	"github.com/Prrost/FinalAP2/order-service/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func RunHTTP(uc *usecase.OrderUsecase, port string) {
	r := gin.Default()

	r.POST("/orders", func(c *gin.Context) {
		var req struct {
			UserEmail string `json:"user_email"`
			BookID    int64  `json:"book_id"`
			DueDays   int    `json:"due_days"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		order, err := uc.CreateOrder(req.UserEmail, req.BookID, req.DueDays)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	})

	r.PUT("/orders/:id/return", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		order, err := uc.ReturnOrder(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	})

	r.GET("/orders/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		order, err := uc.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusOK, order)
	})

	r.GET("/orders/waiting", func(c *gin.Context) {
		bookID, _ := strconv.ParseInt(c.Query("book_id"), 10, 64)
		list, err := uc.ListWaiting(bookID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	})

	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.Run(port)
}
