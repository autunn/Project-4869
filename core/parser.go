package core

import (
	"regexp"
	// 如果你下面没用到 strings.TrimSpace 这种函数，请务必删掉下面这行
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

// 确保用到了 strings，否则报错
func CleanName(name string) string {
	return strings.TrimSpace(name)
}