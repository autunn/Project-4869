package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CD2Config struct {
	Host        string `json:"host"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	TargetDrive string `json:"target_drive"`
}

func AddToCD2(conf CD2Config, magnet string) (string, error) {
	if conf.Host == "" || conf.Username == "" {
		return "", errors.New("CD2配置不完整")
	}
	
	baseURL := strings.TrimRight(conf.Host, "/")
	client := &http.Client{Timeout: 10 * time.Second}

	// 1. Login
	loginPayload := map[string]string{"userName": conf.Username, "password": conf.Password}
	loginBody, _ := json.Marshal(loginPayload)
	resp, err := client.Post(baseURL+"/api/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("登录失败: HTTP %d", resp.StatusCode)
	}

	// 2. Get Drives
	resp, err = client.Get(baseURL + "/api/getDrives")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	var drives []map[string]interface{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &drives); err != nil {
		return "", err
	}

	var driveID string
	var driveName string

	target := conf.TargetDrive
	if target == "" {
		target = "115"
	}

	for _, d := range drives {
		name, _ := d["name"].(string)
		if name == "" {
			name, _ = d["displayName"].(string)
		}
		
		if strings.Contains(name, target) {
			driveID, _ = d["id"].(string)
			driveName = name
			break
		}
	}

	if driveID == "" {
		if len(drives) == 1 {
			driveID, _ = drives[0]["id"].(string)
			driveName, _ = drives[0]["name"].(string)
		} else {
			return "", fmt.Errorf("未找到包含 '%s' 的网盘", target)
		}
	}

	// 3. Add Task
	taskPayload := map[string]string{
		"driveId":      driveID,
		"url":          magnet,
		"parentFileId": "root",
	}
	taskBody, _ := json.Marshal(taskPayload)
	
	resp, err = client.Post(baseURL+"/api/addOfflineFile", "application/json", bytes.NewBuffer(taskBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return fmt.Sprintf("已推送到【%s】离线下载", driveName), nil
	}
	
	bodyBytes, _ = io.ReadAll(resp.Body)
	return "", fmt.Errorf("添加任务失败: %s", string(bodyBytes))
}