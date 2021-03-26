package guide

import (
	"github.com/rs/zerolog"
	"time"
)

type Client struct {
	publicURL string
	deviceID  string

	timeout time.Duration
	log     zerolog.Logger
}

func New(publicURL string, deviceID string, log zerolog.Logger) *Client {
	return &Client{
		publicURL: publicURL,
		timeout:   2 * time.Minute,

		log: log,
	}
}
