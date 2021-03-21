package web

import (
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
	r.GET("/playlist.m3u8", c.Playlist)
	r.GET("/stream.m3u8", c.Stream)
}

func (c *Client) Logger() gin.HandlerFunc {
	return func(g *gin.Context) {
		t := time.Now()
		// before
		g.Next()
		// after
		l := time.Since(t)

		if g.Request == nil {
			return
		}

		// errors
		if len(g.Errors) > 0 {
			errors := make([]error, 0)
			for _, err := range g.Errors {
				errors = append(errors, err.Err)
			}

			c.log.Error().
				Errs("errors", errors).
				Str("url", g.Request.RequestURI).
				Int("status", g.Writer.Status()).
				Str("ip", g.ClientIP()).
				Str("latency", l.String()).
				Msg("Request failed")
			return
		}

		// processed
		c.log.Info().
			Str("url", g.Request.RequestURI).
			Int("status", g.Writer.Status()).
			Str("ip", g.ClientIP()).
			Str("latency", l.String()).
			Msg("Request processed")
	}
}
