package util

import (
	"html"
	"regexp"
	"strings"
)

var (
	reSmoothStreamsChannelName *regexp.Regexp
)

func init() {
	reg1, err := regexp.Compile(`[^\w| |-]`)
	if err != nil {
		panic(err)
	}
	reSmoothStreamsChannelName = reg1
}

func SanitizeChannelName(value string) string {
	v := html.UnescapeString(value)
	return strings.TrimSpace(reSmoothStreamsChannelName.ReplaceAllString(v, ""))
}
