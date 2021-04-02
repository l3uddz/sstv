package web

import (
	"github.com/dustin/go-humanize"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/l3uddz/sstv/logger"
	"github.com/l3uddz/sstv/smoothstreams"
	"github.com/rs/zerolog"
	"golang.org/x/sync/singleflight"
	"strings"
	"time"
)

type Client struct {
	ss         *smoothstreams.Client
	forceProxy bool

	store *persistence.InMemoryStore
	sfg   *singleflight.Group
	log   zerolog.Logger
}

func New(ss *smoothstreams.Client, forceProxy bool) *Client {
	return &Client{
		ss:         ss,
		forceProxy: forceProxy,

		store: persistence.NewInMemoryStore(time.Second),
		sfg:   &singleflight.Group{},
		log:   logger.New(""),
	}
}

func (c *Client) SetHandlers(r *gin.Engine) {
	// core
	r.GET("/playlist.m3u8", c.Playlist)
	r.GET("/stream/:channel", c.WithSanitizedRawQuery(c.Stream))
	r.GET("/epg.xml", cache.CachePage(c.store, time.Minute, c.EPG))
	// plex
	r.GET("/lineup.json", c.Lineup)
	r.GET("/lineup_status.json", c.LineupStatus)
	r.GET("/lineup.post", c.LineupPost)
	r.POST("/lineup.post", c.LineupPost)
	r.GET("/discover.json", c.Discover)
	r.GET("/device.xml", c.Device)
	r.GET("/", c.Device)
}

func (c Client) WithSanitizedRawQuery(next func(*gin.Context)) gin.HandlerFunc {
	return func(g *gin.Context) {
		// plex dvr appends its transcode query param with a ?, even if the url already contains a question mark
		// this causes issues when it comes to decoding, and thus, this crude fix is required
		if strings.Contains(g.Request.URL.RawQuery, "?") {
			g.Request.URL.RawQuery = strings.Replace(g.Request.URL.RawQuery, "?", "&", -1)
		}

		next(g)
	}
}

func (c *Client) Logger() gin.HandlerFunc {
	return func(g *gin.Context) {
		// log request
		rl := c.log.With().
			Str("ip", g.ClientIP()).
			Str("uri", g.Request.RequestURI).
			Logger()

		rl.Debug().Msg("Request received")

		// handle request
		t := time.Now()
		g.Next()
		l := time.Since(t)

		// log errors
		switch {
		case len(g.Errors) > 0:
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

		case g.Writer.Status() >= 400 && g.Writer.Status() <= 599:
			rl.Error().
				Int("status", g.Writer.Status()).
				Str("duration", l.String()).
				Msg("Request failed")
			return
		}

		// log outcome
		rl.Info().
			Str("size", humanize.IBytes(uint64(g.Writer.Size()))).
			Int("status", g.Writer.Status()).
			Str("duration", l.String()).
			Msg("Request processed")
	}
}
