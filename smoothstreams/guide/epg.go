package guide

import (
	"encoding/json"
	"fmt"
	"github.com/l3uddz/sstv/build"
	"github.com/lucperkins/rek"
	"net/http"
	"sort"
	"strconv"
)

const (
	SPORT int = iota + 1
)

func (c *Client) GetEPG(opts *EpgOptions) ([]Channel, error) {
	// determine request url
	requestURL := ""
	switch opts.Type {
	case SPORT:
		// sports epg
		requestURL = "https://fast-guide.smoothstreams.tv/feed.json"
	default:
		// default epg
		requestURL = fmt.Sprintf("https://fast-guide.smoothstreams.tv/altepg/feedall%d.json", opts.Days)
	}

	// create epg request
	resp, err := rek.Get(requestURL, rek.Timeout(c.timeout), rek.UserAgent(build.UserAgent))
	if err != nil {
		return nil, fmt.Errorf("request epg: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("validate epg response: %s", resp.Status())
	}

	// decode epg response
	b := make(map[string]Channel, 0)
	if err := json.NewDecoder(resp.Body()).Decode(&b); err != nil {
		return nil, fmt.Errorf("decode epg response: %w", err)
	}

	// transform epg response
	keys := make([]string, 0)
	for k := range b {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		numA, _ := strconv.Atoi(keys[i])
		numB, _ := strconv.Atoi(keys[j])
		return numA < numB
	})

	channels := make([]Channel, 0)
	for _, k := range keys {
		channels = append(channels, b[k])
	}

	return channels, nil
}
