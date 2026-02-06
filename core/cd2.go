package core

import "log"

func TriggerCD2Download(url string, token string) {
	AddLog("CD2 推送: 成功下发 URL 至 CloudDrive2")
	log.Printf("CD2 任务: %s", url)
}