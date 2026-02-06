package db

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// 确保路径为 data/ 否则 Docker 挂载会找不到
	DB, err = gorm.Open(sqlite.Open("data/p4869.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	log.Println("数据库初始化成功")
}