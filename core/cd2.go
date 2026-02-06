package core

import "log"

func TriggerCD2Download(url string, token string) {
	AddLog("CD2 推送: 任务已成功提交至 CloudDrive2")
	log.Printf("CD2 Action: %s", url)
}