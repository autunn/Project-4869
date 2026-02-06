package core

import (
	"regexp"
	"strings"
)

func ParseEpisode(fileName string) string {
	re := regexp.MustCompile(`\[(\d{1,4})\]`)
	matches := re.FindStringSubmatch(fileName)
	if len(matches) > 1 {
		return matches[1]
	}
	return "unknown"
}

func CleanName(name string) string {
	return strings.TrimSpace(name)
}