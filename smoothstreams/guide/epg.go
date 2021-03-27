package guide

import (
	"encoding/json"
	"fmt"
	"github.com/l3uddz/sstv/build"
	"github.com/lucperkins/rek"
	"sort"
	"strconv"
	"strings"
)

func (c *Client) GetEPG(opts *EpgOptions) ([]Channel, error) {
	// determine request url
	requestURL := fmt.Sprintf("https://fast-guide.smoothstreams.tv/altepg/feedall%d.json", opts.Days)

	if strings.EqualFold(opts.Type, "sport") {
		requestURL = "https://fast-guide.smoothstreams.tv/feed.json"
	}

	// create epg request
	resp, err := rek.Get(requestURL, rek.Timeout(c.timeout), rek.UserAgent(build.UserAgent))
	if err != nil {
		return nil, fmt.Errorf("request standard epg: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate standard epg response: %s", resp.Status())
	}

	// decode epg response
	b := make(map[string]Channel, 0)
	if err := json.NewDecoder(resp.Body()).Decode(&b); err != nil {
		return nil, fmt.Errorf("decode standard epg response: %w", err)
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
