package core

import (
	"regexp"
	"strconv"
	"strings"
)

type ParsedMeta struct {
	Episode    string
	Resolution string
	Container  string
	Subtitle   string
	SourceType string
}

func ParseTitle(title string) ParsedMeta {
	meta := ParsedMeta{}
	if title == "" {
		return meta
	}

	// 1. Episode
	// 匹配剧场版 (Mx)
	reMovie := regexp.MustCompile(`(?i)(?:M|Movie|剧场版)[\s_]*(\d{1,2})`)
	// 匹配普通集数 (123, [123], 【123】)
	reEp := regexp.MustCompile(`(?:^|[\[【\s])(\d{3,4})(?:[\]】\s]|$)`)

	if matches := reMovie.FindStringSubmatch(title); len(matches) > 1 {
		meta.Episode = "M" + matches[1]
	} else if matches := reEp.FindStringSubmatch(title); len(matches) > 1 {
		num, _ := strconv.Atoi(matches[1])
		// 排除年份 (199x - 205x)
		if num < 1990 || num > 2050 {
			meta.Episode = matches[1]
		}
	}

	// 2. Resolution
	reRes := regexp.MustCompile(`(?i)(1080[Pp]|720[Pp]|2160[Pp]|4[Kk])`)
	if matches := reRes.FindStringSubmatch(title); len(matches) > 1 {
		meta.Resolution = strings.ToUpper(matches[1])
	}

	// 3. Container
	reCont := regexp.MustCompile(`(?i)(MKV|MP4|AVI)`)
	if matches := reCont.FindStringSubmatch(title); len(matches) > 1 {
		meta.Container = strings.ToUpper(matches[1])
	}

	// 4. Subtitle
	upperTitle := strings.ToUpper(title)
	if strings.Contains(upperTitle, "CHS_JP") || strings.Contains(upperTitle, "简日") {
		meta.Subtitle = "CHS_JP"
	} else if strings.Contains(upperTitle, "CHT_JP") || strings.Contains(upperTitle, "繁日") {
		meta.Subtitle = "CHT_JP"
	} else if strings.Contains(upperTitle, "CHS") || strings.Contains(upperTitle, "简体") {
		meta.Subtitle = "CHS"
	} else if strings.Contains(upperTitle, "CHT") || strings.Contains(upperTitle, "繁体") {
		meta.Subtitle = "CHT"
	} else if strings.Contains(upperTitle, "JP") {
		meta.Subtitle = "JP"
	}

	// 5. Source
	reSrc := regexp.MustCompile(`(?i)(WEBRIP|HDTV|BDRIP|BLURAY|DVDRIP)`)
	if matches := reSrc.FindStringSubmatch(title); len(matches) > 1 {
		meta.SourceType = strings.ToUpper(matches[1])
	} else if strings.Contains(upperTitle, "WEB-DL") {
		meta.SourceType = "WEBRIP"
	}

	return meta
}