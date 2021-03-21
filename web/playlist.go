package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Client) Playlist(g *gin.Context) {
	// channel requested?
	b := new(struct {
		Channel int `form:"channel"`
		Type    int `form:"type,omitempty"`
	})

	if g.ShouldBindQuery(b) == nil && b.Channel > 0 {
		// get requested channel stream link
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

	// generate playlist
	playlist, err := c.ss.Guide.GeneratePlaylist()
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("generate playlist: %w", err))
	}

	g.Data(http.StatusOK, "application/x-mpegURL; charset=utf-8", []byte(playlist))
}
