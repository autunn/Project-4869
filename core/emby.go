package core

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// EmbyConfig 定义配置结构
type EmbyConfig struct {
	URL    string
	APIKey string
}

// CheckEmby 检查 Emby 连通性
func CheckEmby(cfg EmbyConfig) error {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/System/Info?api_key=%s", cfg.URL, cfg.APIKey), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Emby 返回状态码: %d", resp.StatusCode)
	}

	log.Println("Emby 服务连接正常")
	return nil
}