package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Client) Discover(g *gin.Context) {
	// generate discover
	discover, err := c.ss.Guide.GenerateDiscover()
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("generate discover: %w", err))
		return
	}

	// return discover
	g.Data(http.StatusOK, "application/json; charset=utf-8", []byte(discover))
}
