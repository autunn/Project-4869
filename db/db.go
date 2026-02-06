package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Magnet struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	MagnetLink   string    `gorm:"uniqueIndex:idx_mag_ep" json:"magnet_link"`
	Episode      string    `gorm:"uniqueIndex:idx_mag_ep" json:"episode"`
	EpisodeTitle string    `json:"episode_title"`
	Resolution   string    `json:"resolution"`
	Container    string    `json:"container"`
	Subtitle     string    `json:"subtitle"`
	SourceType   string    `json:"source_type"`
	RawTitle     string    `json:"raw_title"`
	PublishDate  string    `json:"publish_date"`
	CreatedAt    time.Time `json:"created_at"`
}

func InitDB() {
	// 确保数据目录存在
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	var err error
	// 开启 WAL 模式 (Write-Ahead Logging) 以支持高并发读写
	// busy_timeout 设置为 5000ms，防止 locked 错误
	dsn := "data/project4869.db?_journal_mode=WAL&_busy_timeout=5000"
	
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移表结构
	DB.AutoMigrate(&Magnet{})
	log.Println("Database initialized with WAL mode.")
}