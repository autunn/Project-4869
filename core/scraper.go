package core

import (
	"log"
)

func RunScraper() {
	// 确保这里至少用了一次 log，否则就把上面的 import "log" 删掉
	log.Println(">>> 启动自动化任务轮询...")
	MonitorResources()
}