package core

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type EmbyConfig struct {
	URL    string
	APIKey string
}

func CheckEmby(cfg EmbyConfig) error {
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/System/Info?api_key=%s", cfg.URL, cfg.APIKey)
	
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Emby 错误状态码: %d", resp.StatusCode)
	}

	log.Println("Emby 连接正常")
	return nil
}