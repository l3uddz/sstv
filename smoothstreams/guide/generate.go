package guide

import (
	"encoding/json"
	"fmt"
	"github.com/beevik/etree"
	"github.com/l3uddz/sstv"
	"github.com/l3uddz/sstv/smoothstreams/util"
	"strings"
)

type PlaylistOptions struct {
	Type  int  `form:"type,omitempty"`
	Proxy bool `form:"proxy,omitempty"`
}

func (c *Client) GeneratePlaylist(opts *PlaylistOptions) (string, error) {
	// retrieve channels
	channels, err := c.GetChannels()
	if err != nil {
		return "", fmt.Errorf("get channels: %w", err)
	}

	// generate playlist
	data := []string{"#EXTM3U"}
	for _, channel := range channels {
		// prepare channel name
		name := util.SanitizeString(channel.Name)
		if strings.Index(name, " - ") >= 0 {
			name = strings.TrimSpace(name[strings.Index(name, " - ")+3:])
		}

		if name == "" {
			name = fmt.Sprintf("Channel %s", channel.Number)
		}

		// prepare channel logo
		logo := channel.Image
		if !strings.HasSuffix(logo, ".png") {
			logo = "https://i.imgur.com/UyrGfW2.png"
		}

		// add channel to playlist data
		data = append(data, fmt.Sprintf(
			"#EXTINF:-1 tvg-id=%q tvg-name=%q tvg-logo=%q tvg-chno=%q channel-id=%q group-title=%q,%s",
			channel.Number, name, logo, channel.Number, channel.Number, "SmoothStreams", name))
		data = append(data, sstv.JoinURL(c.publicURL,
			fmt.Sprintf("stream.m3u8?channel=%s&type=%d&proxy=%v", channel.Number, opts.Type, opts.Proxy)))
	}

	return strings.Join(data, "\n"), nil
}

func (c *Client) GenerateLineup(opts *PlaylistOptions) (string, error) {
	type lineup struct {
		GuideNumber string
		GuideName   string
		URL         string
	}

	// retrieve channels
	channels, err := c.GetChannels()
	if err != nil {
		return "", fmt.Errorf("get channels: %w", err)
	}

	// generate lineup
	data := make([]lineup, 0)
	for _, channel := range channels {
		// prepare channel name
		name := util.SanitizeString(channel.Name)
		if strings.Index(name, " - ") >= 0 {
			name = strings.TrimSpace(name[strings.Index(name, " - ")+3:])
		}

		if name == "" {
			name = fmt.Sprintf("Channel %s", channel.Number)
		}

		// add channel to lineup
		data = append(data, lineup{
			GuideNumber: channel.Number,
			GuideName:   name,
			URL: sstv.JoinURL(c.publicURL,
				fmt.Sprintf("stream.m3u8?channel=%s&type=%d&plex=1", channel.Number, opts.Type)),
		})
	}

	// marshal
	b, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshal lineup: %w", err)
	}

	return string(b), nil
}

func (c *Client) GenerateLineupStatus() (string, error) {
	type lineupStatus struct {
		ScanInProgress int
		ScanPossible   int
		Source         string
		SourceList     []string
	}

	// generate lineup status
	data := &lineupStatus{
		ScanInProgress: 0,
		ScanPossible:   1,
		Source:         "Cable",
		SourceList:     []string{"Cable"},
	}

	// marshal
	b, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshal lineup_status: %w", err)
	}

	return string(b), nil
}

type EpgOptions struct {
	Days int `form:"days,omitempty"`
}

func (c *Client) GenerateEPG(opts *EpgOptions) (string, error) {
	// option defaults
	if opts.Days == 0 || opts.Days > 5 {
		// 5 is the maximum days
		opts.Days = 5
	}

	// retrieve channels (with epg data)
	channels, err := c.GetEPG(opts.Days)
	if err != nil {
		return "", fmt.Errorf("get epg: %w", err)
	}

	// prepare generate epg
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	tvd := doc.CreateElement("tv")
	tvd.CreateAttr("date", util.CurrentXmlTvTime())
	tvd.CreateAttr("generator_info_name", "sstv")

	// generate epg
	for _, channel := range channels {
		// prepare channel name
		name := util.SanitizeString(channel.Name)
		if strings.Index(name, " - ") >= 0 {
			name = strings.TrimSpace(name[strings.Index(name, " - ")+3:])
		}

		if name == "" {
			name = fmt.Sprintf("Channel %s", channel.Number)
		}

		// prepare channel logo
		logo := channel.Image
		if !strings.HasSuffix(logo, ".png") {
			logo = "https://i.imgur.com/UyrGfW2.png"
		}

		// create channel element
		chd := tvd.CreateElement("channel")
		chd.CreateAttr("id", channel.Number)
		chd.CreateElement("display-name").CreateText(name)
		chd.CreateElement("icon").CreateAttr("src", logo)

		// create programme elements
		for _, programme := range channel.Programmes {
			// prepare programme element
			name = util.SanitizeString(programme.Name)
			start, err := util.TimeStringToXmlTvTime(programme.StartTime)
			if err != nil {
				continue
			}
			end, err := util.TimeStringToXmlTvTime(programme.EndTime)
			if err != nil {
				continue
			}

			// create programme element
			pgd := tvd.CreateElement("programme")
			pgd.CreateAttr("channel", programme.Channel)
			pgd.CreateAttr("start", start)
			pgd.CreateAttr("stop", end)

			pgd.CreateElement("title").CreateText(name)

			if programme.Description != "" {
				pgd.CreateElement("desc").CreateText(programme.Description)
			}

			if programme.Category != "" {
				pgd.CreateElement("category").CreateText(programme.Category)
			}

			// create episode element
			epd := pgd.CreateElement("episode-num")
			epd.CreateAttr("system", "dd_progid")
			epd.CreateText(programme.Id)
		}
	}

	// return epg
	doc.Indent(2)
	data, err := doc.WriteToString()
	if err != nil {
		return "", fmt.Errorf("write to string: %w", err)
	}

	return data, nil
}
