package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type EmbyConfig struct {
	Host       string `json:"host"`
	ApiKey     string `json:"api_key"`
	TmdbId     string `json:"tmdb_id"`
	MaxEpisode int    `json:"max_episode"`
}

type EmbyResponse struct {
	Items []struct {
		Id          string `json:"Id"`
		IndexNumber int    `json:"IndexNumber,omitempty"`
	} `json:"Items"`
}

func CheckEmbyMissing(conf EmbyConfig) ([]int, int, error) {
	if conf.Host == "" {
		return nil, 0, errors.New("Emby Host未配置")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	headers := http.Header{}
	headers.Set("X-Emby-Token", conf.ApiKey)

	// 1. Find Series ID by TMDB ID
	url := fmt.Sprintf("%s/Items?Recursive=true&IncludeItemTypes=Series&AnyProviderIdEquals=tmdb.%s", strings.TrimRight(conf.Host, "/"), conf.TmdbId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header = headers
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var seriesResp EmbyResponse
	if err := json.NewDecoder(resp.Body).Decode(&seriesResp); err != nil {
		return nil, 0, err
	}
	if len(seriesResp.Items) == 0 {
		return nil, 0, fmt.Errorf("Emby中未找到该剧集 (TMDB: %s)", conf.TmdbId)
	}
	seriesID := seriesResp.Items[0].Id

	// 2. Get All Episodes
	epUrl := fmt.Sprintf("%s/Items?ParentId=%s&Recursive=true&IncludeItemTypes=Episode&Fields=IndexNumber", strings.TrimRight(conf.Host, "/"), seriesID)
	req, _ = http.NewRequest("GET", epUrl, nil)
	req.Header = headers
	
	resp, err = client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var epsResp EmbyResponse
	if err := json.NewDecoder(resp.Body).Decode(&epsResp); err != nil {
		return nil, 0, err
	}

	existingMap := make(map[int]bool)
	for _, item := range epsResp.Items {
		existingMap[item.IndexNumber] = true
	}

	// 3. Calculate Missing
	var missing []int
	max := conf.MaxEpisode
	if max == 0 {
		max = 1200 // Default fallback
	}

	for i := 1; i <= max; i++ {
		if !existingMap[i] {
			missing = append(missing, i)
		}
	}
	
	sort.Ints(missing)
	return missing, len(existingMap), nil
}

import "strings"