package util

import (
	"html"
	"regexp"
	"strings"
)

var (
	reSmoothStreamsString *regexp.Regexp
)

func init() {
	reg1, err := regexp.Compile(`[^\w| |-]`)
	if err != nil {
		panic(err)
	}
	reSmoothStreamsString = reg1
}

func SanitizeString(value string) string {
	v := html.UnescapeString(value)
	return strings.TrimSpace(reSmoothStreamsString.ReplaceAllString(v, ""))
}
