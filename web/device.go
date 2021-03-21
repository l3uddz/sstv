package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Client) Device(g *gin.Context) {
	// generate device
	device, err := c.ss.Guide.GenerateDevice()
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("generate device: %w", err))
	}

	// return device
	g.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(device))
}
