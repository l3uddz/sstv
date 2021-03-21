package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Client) Playlist(g *gin.Context) {
	// parse query
	b := new(struct {
		Type int `form:"type,omitempty"`
	})

	if err := g.ShouldBindQuery(b); err != nil {
		g.AbortWithError(http.StatusBadRequest, fmt.Errorf("bind query: %w", err))
		return
	}

	// generate playlist
	playlist, err := c.ss.Guide.GeneratePlaylist(b.Type)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("generate playlist: %w", err))
	}

	// return playlist
	g.Data(http.StatusOK, "application/x-mpegURL; charset=utf-8", []byte(playlist))
}
