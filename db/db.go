package db

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type SystemConfig struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	CD2Token string `json:"cd2_token"`
	EmbyURL  string `json:"emby_url"`
}

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/p4869.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	DB.AutoMigrate(&SystemConfig{})
}