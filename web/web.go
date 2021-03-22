package web

import (
	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"github.com/l3uddz/sstv/logger"
	"github.com/l3uddz/sstv/smoothstreams"
	"github.com/rs/zerolog"
	"time"
)

type Client struct {
	ss *smoothstreams.Client

	log zerolog.Logger
}

func New(ss *smoothstreams.Client) *Client {
	return &Client{
		ss:  ss,
		log: logger.New(""),
	}
}

func (c *Client) SetHandlers(r *gin.Engine) {
	// core
	r.GET("/playlist.m3u8", c.Playlist)
	r.GET("/stream.m3u8", c.Stream)
	r.GET("/epg.xml", c.EPG)
	// plex
	r.GET("/lineup.json", c.Lineup)
	r.GET("/lineup_status.json", c.LineupStatus)
	r.GET("/lineup.post", c.LineupPost)
	r.POST("/lineup.post", c.LineupPost)
	r.GET("/discover.json", c.Discover)
	r.GET("/device.xml", c.Device)
	r.GET("/", c.Device)
}

func (c *Client) Logger() gin.HandlerFunc {
	return func(g *gin.Context) {
		if g.Request == nil {
			g.Next()
			return
		}

		rl := c.log.With().
			Str("ip", g.ClientIP()).
			Str("uri", g.Request.RequestURI).
			Logger()

		rl.Debug().Msg("Request received")

		t := time.Now()
		g.Next()
		l := time.Since(t)

		if len(g.Errors) > 0 {
			errors := make([]error, 0)
			for _, err := range g.Errors {
				errors = append(errors, err.Err)
			}

			rl.Error().
				Errs("errors", errors).
				Int("status", g.Writer.Status()).
				Str("duration", l.String()).
				Msg("Request failed")
			return
		}

		rl.Info().
			Str("size", humanize.IBytes(uint64(g.Writer.Size()))).
			Int("status", g.Writer.Status()).
			Str("duration", l.String()).
			Msg("Request processed")
	}
}
