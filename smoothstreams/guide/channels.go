package guide

import (
	"encoding/json"
	"fmt"
	"github.com/l3uddz/sstv/build"
	"github.com/lucperkins/rek"
	"sort"
	"strconv"
)

type Channel struct {
	Number     string      `json:"channel_id"`
	Name       string      `json:"name"`
	Image      string      `json:"img"`
	Programmes []Programme `json:"items,omitempty"`
}

type Programme struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	StartTime   string `json:"time"`
	EndTime     string `json:"end_time"`
	Channel     string `json:"channel"`
	Category    string `json:"category"`
}

func (c *Client) GetChannels() ([]Channel, error) {
	// create guide request
	resp, err := rek.Get("https://fast-guide.smoothstreams.tv/feed.json", rek.Timeout(c.timeout),
		rek.UserAgent(build.UserAgent))
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
