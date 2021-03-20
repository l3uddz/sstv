package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Client) Playlist(g *gin.Context) {
	// decode query
	b := new(struct {
		Channel int `form:"channel"`
		Type    int `form:"type,omitempty"`
	})

	// channel requested?
	if g.ShouldBindQuery(b) == nil {
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

	// send channel playlist
}
