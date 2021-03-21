package guide

import (
	"encoding/json"
	"fmt"
	"github.com/lucperkins/rek"
	"sort"
	"strconv"
)

type Channel struct {
	Number string `json:"channel_id"`
	Name   string `json:"name"`
	Image  string `json:"img"`
}

func (c *Client) GetChannels() ([]Channel, error) {
	// create guide request
	resp, err := rek.Get("https://fast-guide.smoothstreams.tv/feed.json", rek.Timeout(c.timeout))
	if err != nil {
		return nil, fmt.Errorf("request guide: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate guide response: %s", resp.Status())
	}

	// decode guide response
	b := make(map[string]Channel, 0)
	if err := json.NewDecoder(resp.Body()).Decode(&b); err != nil {
		return nil, fmt.Errorf("decode guide response: %w", err)
	}

	// transform guide response
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
