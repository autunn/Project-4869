package core

import (
	"regexp"
	"strings"
)

// ParseEpisode 从文件名解析集数
func ParseEpisode(fileName string) string {
	// 简单的正则匹配示例：第1000集
	re := regexp.MustCompile(`\[(\d{1,4})\]`)
	matches := re.FindStringSubmatch(fileName)
	if len(matches) > 1 {
		return matches[1]
	}
	return "unknown"
}

// CleanName 格式化名称
func CleanName(name string) string {
	return strings.TrimSpace(name)
}