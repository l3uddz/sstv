package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Client) Playlist(g *gin.Context) {
	// generate playlist
	playlist, err := c.ss.Guide.GeneratePlaylist()
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("generate playlist: %w", err))
	}

	// return playlist
	g.Data(http.StatusOK, "application/x-mpegURL; charset=utf-8", []byte(playlist))
}
