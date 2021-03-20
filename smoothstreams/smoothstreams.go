package smoothstreams

import (
	"fmt"
	"github.com/l3uddz/sstv/logger"
	"github.com/l3uddz/sstv/smoothstreams/token"
	"github.com/rs/zerolog"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Site     string `yaml:"site"`
	Server   string `yaml:"server"`

	Verbosity string `yaml:"verbosity"`
}

type Client struct {
	token *token.Client

	log zerolog.Logger
}

func New(c Config) (*Client, error) {
	l := logger.New(c.Verbosity)

	// token
	t, err := token.New(c.Username, c.Password, c.Site, l)
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}

	return &Client{
		token: t,
		log:   l,
	}, nil
}
