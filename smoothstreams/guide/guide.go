package guide

import (
	"github.com/rs/zerolog"
	"time"
)

type Client struct {
	publicURL string
	timeout   time.Duration

	log zerolog.Logger
}

func New(publicURL string, log zerolog.Logger) *Client {
	return &Client{
		publicURL: publicURL,
		timeout:   2 * time.Minute,

		log: log,
	}
}
