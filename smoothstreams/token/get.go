package token

import (
	"encoding/json"
	"fmt"
	"github.com/l3uddz/sstv"
	"github.com/l3uddz/sstv/build"
	"github.com/lucperkins/rek"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) Get() (string, error) {
	// acquire lock
	c.mtx.Lock()
	defer c.mtx.Unlock()

	// check stored token
	if c.hash != "" && time.Now().UTC().Before(c.expiry.UTC()) {
		return c.hash, nil
	}

	// generate token request url
	tokenURL, err := sstv.URLWithQuery(c.authURL(), url.Values{
		"username": []string{c.Username},
		"password": []string{c.Password},
		"site":     []string{c.Site},
	})
	if err != nil {
		return c.hash, fmt.Errorf("token url: %w", err)
	}

	c.log.Trace().
		Str("url", tokenURL).
		Msg("Requesting token")

	// create token request
	resp, err := rek.Get(tokenURL, rek.Timeout(c.timeout), rek.UserAgent(build.UserAgent))
	if err != nil {
		return c.hash, fmt.Errorf("request token: %w", err)
	}
	defer resp.Body().Close()

	// validate token response
	if resp.StatusCode() != http.StatusOK {
		return c.hash, fmt.Errorf("validate token response: %s", resp.Status())
	}

	// decode token response
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

	// update stored token
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
