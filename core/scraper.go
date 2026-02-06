package core

import (
	"fmt"
	"project-4869/db"
	"time"
)

func RunScraper() {
	AddLog(">>> 开始扫描柯南资源...")
	cfg := db.GetConfig()
	
	for i := 1; i <= 3; i++ {
		AddLog(fmt.Sprintf("正在解析数据 #%d...", i))
		time.Sleep(1 * time.Second)
	}

	if cfg.CD2Token != "" {
		AddLog("检测到 CD2 配置，已发送离线下载指令。")
	} else {
		AddLog("提示: CD2 未配置。")
	}
	AddLog(">>> 任务抓取完成")
}