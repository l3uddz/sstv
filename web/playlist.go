package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/l3uddz/sstv/smoothstreams/guide"
	"net/http"
)

func (c *Client) Playlist(g *gin.Context) {
	// parse query
	b := new(guide.PlaylistOptions)

	if err := g.ShouldBindQuery(b); err != nil {
		g.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind query: %w", err))
		return
	}

	// generate playlist
	playlist, err := c.ss.Guide.GeneratePlaylist(b)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("generate playlist: %w", err))
		return
	}

	// return playlist
	g.Data(http.StatusOK, "application/x-mpegURL; charset=utf-8", []byte(playlist))
}
