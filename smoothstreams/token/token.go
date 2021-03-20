package token

import (
	"fmt"
	"github.com/rs/zerolog"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Username string
	Password string
	Site     string

	hash   string
	code   string
	expiry time.Time

	mtx     *sync.Mutex
	timeout time.Duration
	log     zerolog.Logger
}

func New(user string, pass string, site string, log zerolog.Logger) (*Client, error) {
	c := &Client{
		Username: user,
		Password: pass,
		Site:     site,

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
	if strings.Contains(c.Site, "mma") {
		return "https://www.mma-tv.net/loginForm.php"
	}
	return "https://auth.SmoothStreams.tv/hash_api.php"
}
