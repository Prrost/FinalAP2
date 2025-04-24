package proxi

import (
	"api-gateway/gateway/Response"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

func ProxiRequest(targetURL string, c *gin.Context) {
	const op = "ProxiRequest"

	fullUrl, err := url.JoinPath(targetURL, c.Request.URL.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response.Err{Error: fmt.Sprintf("%s: %s", op, err.Error())})
		return
	}

	if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
		fullUrl += "?" + rawQuery
	}

	req, err := http.NewRequest(c.Request.Method, fullUrl, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response.Err{Error: fmt.Sprintf("%s: %s", op, err.Error())})
		return
	}

	for name, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response.Err{Error: fmt.Sprintf("%s: %s", op, err.Error())})
		slog.Error("Unable to connect to server. Error:", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response.Err{Error: fmt.Sprintf("%s: %s", op, err.Error())})
		return
	}

	c.Writer.WriteHeader(resp.StatusCode)
	_, err = c.Writer.Write(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response.Err{Error: fmt.Sprintf("%s: %s", op, err.Error())})
	}
}
