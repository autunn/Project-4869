package core

import (
	"log"
)

func TriggerCD2Download(url string, token string) {
	log.Printf("准备通知 CD2 下载资源，URL: %s", url)
	// 具体的 http 请求逻辑暂未启用，因此不引用 net/http
}