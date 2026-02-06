package core

import (
	"log"
)

func RunScraper() {
	log.Println(">>> 启动自动化任务轮询...")
	MonitorResources()
	log.Println(">>> 抓取任务执行完毕")
}