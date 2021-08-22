package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/l3uddz/sstv/build"
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
		Channel int `uri:"channel" binding:"required"`
	})

	if err := g.ShouldBindUri(b); err != nil {
		_ = g.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind uri: %w", err))
		return
	}

	if err := g.ShouldBindQuery(b); err != nil {
		_ = g.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind query: %w", err))
		return
	}

	// validate request
	if b.Channel == 0 {
		_ = g.AbortWithError(http.StatusBadRequest, errors.New("channel was not parsed"))
		return
	}

	// adjust request
	if c.forceProxy && !b.Plex {
		b.Type = stream.MPEG2TS
		b.Proxy = true
	}

	// get stream link
	cl, err := c.ss.Stream.GetLink(b.Channel, b.Server, b.Type)
	if err != nil {
		_ = g.AbortWithError(http.StatusInternalServerError,
			fmt.Errorf("get stream link: %d.%d: %w", b.Channel, b.Type, err))
		return
	}

	// redirect to stream link
	if !b.Proxy {
		g.Redirect(http.StatusTemporaryRedirect, cl)
		return
	} else if b.Type != stream.MPEG2TS {
		// we can only proxy MPEG2TS streams
		_ = g.AbortWithError(http.StatusUnsupportedMediaType, errors.New("stream type cannot be proxied"))
		return
	}

	// proxy stream link
	resp, err := rek.Get(cl, rek.UserAgent(build.UserAgent))
	if err != nil {
		_ = g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("get proxy stream: %w", err))
		return
	}
	defer resp.Body().Close()

	if resp.StatusCode() != http.StatusOK {
		_ = g.AbortWithError(http.StatusServiceUnavailable, fmt.Errorf("get proxy stream: %s", resp.Status()))
		return
	}

	g.Writer.Header().Set("Content-Type", resp.ContentType())
	_, _ = io.Copy(g.Writer, resp.Body())
}
