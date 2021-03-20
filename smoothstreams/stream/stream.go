package stream

import (
	"github.com/l3uddz/sstv/smoothstreams/token"
	"github.com/rs/zerolog"
)

type Client struct {
	Site   string
	Server string

	token *token.Client
	log   zerolog.Logger
}

func New(site string, server string, token *token.Client, log zerolog.Logger) *Client {
	return &Client{
		Site:   site,
		Server: server,

		token: token,
		log:   log,
	}
}
