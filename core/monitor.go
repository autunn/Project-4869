package core

import (
	"log"
	"project-4869/db"
)

func MonitorResources() {
	log.Println("正在查询数据库资源状态...")
	if db.DB != nil {
		var count int64
		db.DB.Table("resources").Count(&count)
		log.Printf("监控到数据库内有 %d 条记录", count)
	}
}