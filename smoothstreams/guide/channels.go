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
	type channel struct {
		Number string `json:"channum"`
		Name   string `json:"channame"`
		Image  string `json:"icon"`
	}

	// create channels request
	resp, err := rek.Get("https://fast-guide.smoothstreams.tv/altepg/channels.json", rek.Timeout(c.timeout),
		rek.UserAgent(build.UserAgent))
	if err != nil {
		return nil, fmt.Errorf("request channels: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("validate channels response: %s", resp.Status())
	}

	// decode channels response
	b := make(map[string]channel, 0)
	if err := json.NewDecoder(resp.Body()).Decode(&b); err != nil {
		return nil, fmt.Errorf("decode channels response: %w", err)
	}

	// transform channels response
	channels := make([]Channel, 0)
	for _, v := range b {
		channels = append(channels, Channel{
			Number: v.Number,
			Name:   v.Name,
			Image:  v.Image,
		})
	}

	sort.Slice(channels, func(i, j int) bool {
		numA, _ := strconv.Atoi(channels[i].Number)
		numB, _ := strconv.Atoi(channels[j].Number)
		return numA < numB
	})

	return channels, nil
}
