package core

import (
	"log"
	"project-4869/db"
)

// MonitorResources 监控数据库中的资源状态
func MonitorResources() {
	log.Println("开始扫描本地资源库...")
	// 示例：查询数据库
	var count int64
	db.DB.Table("resources").Count(&count)
	log.Printf("当前监控中的资源总数: %d", count)
}