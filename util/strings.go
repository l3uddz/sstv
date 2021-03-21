package util

import (
	"html"
	"regexp"
	"strings"
)

var (
	reSmoothStreamsChannelname *regexp.Regexp
)

func init() {
	reg1, err := regexp.Compile(`[^\w| |-]`)
	if err != nil {
		panic(err)
	}
	reSmoothStreamsChannelname = reg1
}

func SanitizeChannelName(value string) string {
	v := html.UnescapeString(value)
	return strings.TrimSpace(reSmoothStreamsChannelname.ReplaceAllString(v, ""))
}
