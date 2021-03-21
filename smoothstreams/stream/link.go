package stream

import "fmt"

const (
	HLS int = iota
	MPEG2TS
	RTMP
)

func (c *Client) GetLink(channel int, streamType int) (string, error) {
	// get current token
	t, err := c.token.Get()
	if err != nil {
		return "", fmt.Errorf("get token: %w", err)
	}

	// generate link
	scheme := "https"
	port := "443"
	playlist := "playlist.m3u8"

	switch streamType {
	case HLS:
		// defaults are already set
	case RTMP:
		scheme = "rtmp"
		port = "3625"
		return fmt.Sprintf("%s://%s.SmoothStreams.tv:%s/%s?wmsAuthSign=%s/ch%sq1.stream",
			scheme, c.Server, port, c.Site, t, fmt.Sprintf("%02d", channel)), nil
	case MPEG2TS:
		playlist = "mpeg.2ts"
	}

	return fmt.Sprintf("%s://%s.SmoothStreams.tv:%s/%s/ch%sq1.stream/%s?wmsAuthSign=%s",
		scheme, c.Server, port, c.Site, fmt.Sprintf("%02d", channel), playlist, t), nil
}
