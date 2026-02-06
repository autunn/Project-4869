package core

import (
	"log"
)

func RunScraper() {
	log.Println(">>> 启动自动化任务轮询...")
	
	// 1. 运行监控
	MonitorResources()
	
	// 2. 这里可以添加具体的抓取逻辑
	log.Println(">>> 任务轮询结束")
}