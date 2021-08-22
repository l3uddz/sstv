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
		_ = g.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind query: %w", err))
		return
	}

	// generate epg (singleflight)
	v, err, _ := c.sfg.Do(g.Request.RequestURI, func() (interface{}, error) {
		return c.ss.Guide.GenerateEPG(b)
	})
	if err != nil {
		_ = g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("generate epg: %w", err))
		return
	}

	// typecast result
	epg, ok := v.(string)
	if !ok {
		_ = g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("typecast epg result"))
		return
	}

	// return epg
	g.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(epg))
}
