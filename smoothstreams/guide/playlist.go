package guide

import (
	"fmt"
	"github.com/l3uddz/sstv"
	"github.com/l3uddz/sstv/util"
	"strings"
)

func (c *Client) GeneratePlaylist() (string, error) {
	// retrieve channels
	channels, err := c.GetChannels()
	if err != nil {
		return "", fmt.Errorf("get channels: %w", err)
	}

	// generate playlist
	data := []string{"#EXTM3U"}
	for _, channel := range channels {
		// prepare channel name
		name := util.SanitizeChannelName(channel.Name)
		if strings.Index(name, " - ") >= 0 {
			name = strings.TrimSpace(name[strings.Index(name, " - ")+3:])
		}

		if name == "" {
			name = "Unknown"
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
		data = append(data, sstv.JoinURL(c.publicURL, fmt.Sprintf("playlist.m3u8?channel=%s", channel.Number)))
	}

	return strings.Join(data, "\n"), nil
}
