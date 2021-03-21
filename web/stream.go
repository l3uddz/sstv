package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Client) Stream(g *gin.Context) {
	// parse query
	b := new(struct {
		Channel int `form:"channel"`
		Type    int `form:"type,omitempty"`
	})

	if err := g.ShouldBindQuery(b); err != nil {
		g.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind query: %w", err))
		return
	}

	// validate query
	if b.Channel == 0 {
		g.AbortWithError(http.StatusBadRequest, errors.New("channel was not parsed"))
		return
	}

	// get stream link
	cl, err := c.ss.Stream.GetLink(b.Channel, b.Type)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError,
			fmt.Errorf("get stream link: %d.%d: %w", b.Channel, b.Type, err))
		return
	}

	// redirect to stream link
	g.Redirect(http.StatusTemporaryRedirect, cl)
	return
}
