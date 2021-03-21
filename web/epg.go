package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/l3uddz/sstv/smoothstreams/guide"
	"net/http"
)

func (c *Client) EPG(g *gin.Context) {
	// parse query
	b := new(guide.EpgOptions)

	if err := g.ShouldBindQuery(b); err != nil {
		g.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind query: %w", err))
		return
	}

	// generate epg
	epg, err := c.ss.Guide.GenerateEPG(b)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("generate epg: %w", err))
		return
	}

	// return epg
	g.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(epg))
}
