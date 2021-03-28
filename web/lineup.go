package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/l3uddz/sstv/smoothstreams/guide"
	"github.com/l3uddz/sstv/smoothstreams/stream"
	"net/http"
)

func (c *Client) Lineup(g *gin.Context) {
	// prepare playlist options
	b := &guide.PlaylistOptions{
		Type: stream.MPEG2TS,
	}

	// generate lineup
	lineup, err := c.ss.Guide.GenerateLineup(b)
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("generate lineup: %w", err))
	}

	// return lineup
	g.Data(http.StatusOK, "application/json; charset=utf-8", []byte(lineup))
}

func (c *Client) LineupStatus(g *gin.Context) {
	// generate lineup_status
	lineupStatus, err := c.ss.Guide.GenerateLineupStatus()
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, fmt.Errorf("generate lineup_status: %w", err))
	}

	// return lineup_status
	g.Data(http.StatusOK, "application/json; charset=utf-8", []byte(lineupStatus))
}

func (c *Client) LineupPost(g *gin.Context) {
	g.String(http.StatusOK, "")
}
