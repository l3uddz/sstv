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
	Number     string      `json:"channum"`
	Name       string      `json:"channame"`
	Image      string      `json:"icon"`
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
	resp, err := rek.Get("https://fast-guide.smoothstreams.tv/altepg/channels.json", rek.Timeout(c.timeout),
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
	channels := make([]Channel, 0)
	for k, _ := range b {
		channels = append(channels, b[k])
	}

	sort.Slice(channels, func(i, j int) bool {
		numA, _ := strconv.Atoi(channels[i].Number)
		numB, _ := strconv.Atoi(channels[j].Number)
		return numA < numB
	})

	return channels, nil
}
