package token

import (
	"encoding/json"
	"fmt"
	"github.com/l3uddz/sstv"
	"github.com/lucperkins/rek"
	"net/url"
	"time"
)

func (c *Client) Get() (string, error) {
	// acquire lock
	c.mtx.Lock()
	defer c.mtx.Unlock()

	// existing token still valid?
	if c.hash != "" && time.Now().UTC().Before(c.expiry.UTC()) {
		c.log.Trace().
			Time("expires", c.expiry).
			Msg("Re-using existing token")
		return c.hash, nil
	}

	// get token url
	tokenUrl, err := sstv.URLWithQuery(c.authUrl(), url.Values{
		"username": []string{c.Username},
		"password": []string{c.Password},
		"site":     []string{c.Site},
	})
	if err != nil {
		return c.hash, fmt.Errorf("token url: %w", err)
	}

	c.log.Trace().
		Str("url", tokenUrl).
		Msg("Requesting token")

	// create token request
	resp, err := rek.Get(tokenUrl, rek.Timeout(c.timeout))
	if err != nil {
		return c.hash, fmt.Errorf("request token: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return c.hash, fmt.Errorf("validate token response: %s", resp.Status())
	}

	// decode response
	b := new(struct {
		Hash  string `json:"hash"`
		Valid int    `json:"valid"`
		Code  string `json:"code,omitempty"`
	})
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return c.hash, fmt.Errorf("decode token response: %w", err)
	}

	if b.Hash == "" {
		return c.hash, fmt.Errorf("validate token response hash: %+v", *b)
	}

	// update token
	c.hash = b.Hash
	c.code = b.Code
	if b.Valid >= 10 {
		c.expiry = time.Now().UTC().Add(time.Duration(b.Valid-5) * time.Minute)
	} else {
		c.expiry = time.Now().UTC().Add(time.Duration(b.Valid) * time.Minute)
	}

	c.log.Debug().
		Time("expires", c.expiry).
		Msg("Retrieved token")
	return c.hash, nil
}
