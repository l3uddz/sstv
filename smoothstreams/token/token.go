package token

import (
	"fmt"
	"github.com/rs/zerolog"
	"sync"
	"time"
)

type Client struct {
	Username string
	Password string
	Site     string
	TokenURL string

	hash   string
	code   string
	expiry time.Time

	mtx     *sync.Mutex
	timeout time.Duration
	log     zerolog.Logger
}

func New(user string, pass string, site string, tokenURL string, log zerolog.Logger) (*Client, error) {
	if tokenURL == "" {
		tokenURL = "https://auth.smoothstreams.tv/hash_api.php"
	}

	c := &Client{
		Username: user,
		Password: pass,
		Site:     site,
		TokenURL: tokenURL,

		mtx:     &sync.Mutex{},
		timeout: 2 * time.Minute,
		log:     log,
	}

	// retrieve initial token
	if _, err := c.Get(); err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	return c, nil
}

func (c *Client) authUrl() string {
	return c.TokenURL
}
