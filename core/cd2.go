package core

import (
	"log"
)

func TriggerCD2(url string, token string) {
	AddLog("准备推送至 CD2...")
	log.Printf("推送 URL: %s, Token: %s", url, token)
}