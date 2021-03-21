package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/l3uddz/sstv/smoothstreams/guide"
	"github.com/l3uddz/sstv/smoothstreams/stream"
	"github.com/lucperkins/rek"
	"io"
	"net/http"
)

func (c *Client) Stream(g *gin.Context) {
	// parse query
	b := new(struct {
		guide.PlaylistOptions
		Channel int `form:"channel"`
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
	if !b.Proxy {
		g.Redirect(http.StatusTemporaryRedirect, cl)
		return
	} else if b.Type != stream.MPEG2TS {
		// we can only proxy MPEG2TS streams
		g.AbortWithError(http.StatusUnsupportedMediaType, errors.New("stream type cannot be proxied"))
		return
	}

	// proxy stream link
	resp, err := rek.Get(cl)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("get proxy stream: %w", err))
		return
	}

	if resp.StatusCode() != http.StatusOK {
		g.AbortWithError(http.StatusServiceUnavailable, fmt.Errorf("get proxy stream: %s", resp.Status()))
		return
	}

	g.Writer.Header().Set("Content-Type", resp.ContentType())
	_, _ = io.Copy(g.Writer, resp.Body())
}
